package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	h "github.com/ThomasITU/DISYSPrep/helpermethod"

	"google.golang.org/grpc"

	"github.com/ThomasITU/DISYSPrep/Proto"
)

const (
	SERVER_PORT     = 5000
	SERVER_LOG_FILE = "serverLog"
)

type Server struct {
	Proto.UnimplementedProtoServiceServer
	port             int
	latestValue      Value
	connectedUsers   []int
	arbiter          sync.Mutex
	lamportTimeStamp int64
}

type Value struct {
	value  int64
	userId int64
}

func main() {

	//init
	initValue := Value{value: -1, userId: -1}
	users := make([]int, 0)
	server := Server{port: SERVER_PORT, latestValue: initValue, connectedUsers: users, arbiter: sync.Mutex{}, lamportTimeStamp: 0}

	fmt.Println("server running")
	listen(&server)
	fmt.Println("main has ended")
}

// JoinService grpc method logic
func (s *Server) JoinService(ctx context.Context, request *Proto.JoinRequest) (*Proto.Response, error) {
	s.arbiter.Lock()
	var msg string
	userId := int(request.GetUserId())
	s.IncrementLamportTimestamp(request.GetTimestamp())
	h.LoggerWithTimestamp("request Received by: "+strconv.Itoa(userId), SERVER_LOG_FILE, s.lamportTimeStamp)
	// check if user id already exist in the array
	for _, user := range s.connectedUsers {
		if user == userId {
			msg = fmt.Sprintf("A user with id: %v has already joined", userId)
			break
		}
	}

	//add userid to slice
	if msg == "" {
		s.connectedUsers = append(s.connectedUsers, userId)
		msg = fmt.Sprintf("Welcome user: %v", userId)
	}
	s.lamportTimeStamp = (s.lamportTimeStamp + 1)
	h.LoggerWithTimestamp(msg, SERVER_LOG_FILE, s.lamportTimeStamp)
	s.arbiter.Unlock()
	return &Proto.Response{Msg: msg, Timestamp: s.lamportTimeStamp}, nil
}

// get value grpc method logic
func (s *Server) GetValue(ctx context.Context, request *Proto.GetRequest) (*Proto.Value, error) {
	s.IncrementLamportTimestamp(request.GetTimestamp())
	h.LoggerWithTimestamp("getvalue request received", SERVER_LOG_FILE, s.lamportTimeStamp)

	s.lamportTimeStamp = (s.lamportTimeStamp + 1)
	value := Proto.Value{CurrentValue: s.latestValue.value, UserId: s.latestValue.userId, Timestamp: s.lamportTimeStamp}

	msg := fmt.Sprintf("value: %v by user: %v", value.GetCurrentValue(), value.GetUserId())
	h.LoggerWithTimestamp("send - "+msg, SERVER_LOG_FILE, s.lamportTimeStamp)
	return &value, nil
}

// set value grpc method logic
func (s *Server) SetValue(ctx context.Context, request *Proto.SetRequest) (*Proto.Response, error) {
	s.arbiter.Lock()
	s.IncrementLamportTimestamp(request.GetTimestamp())
	h.LoggerWithTimestamp("setvalue request received - ", SERVER_LOG_FILE, s.lamportTimeStamp)
	temp := s.latestValue
	s.latestValue = Value{value: request.GetRequestedValue(), userId: request.GetUserId()}

	msg := fmt.Sprintf("Updated the value: %v by user %v to %v by user %v", temp.value, temp.userId, s.latestValue.value, s.latestValue.userId)
	s.lamportTimeStamp = (s.lamportTimeStamp + 1)
	h.LoggerWithTimestamp("send - "+msg, SERVER_LOG_FILE, s.lamportTimeStamp)
	s.arbiter.Unlock()
	return &Proto.Response{Msg: msg, Timestamp: s.lamportTimeStamp}, nil
}

// start server service
func listen(s *Server) {

	//listen on port
	lis, err := net.Listen("tcp", "localhost:"+strconv.Itoa(s.port))
	h.CheckError(err, "server setup net.listen")
	defer lis.Close()

	// register server this is a blocking call
	grpcServer := grpc.NewServer()
	Proto.RegisterProtoServiceServer(grpcServer, s)
	errorMsg := grpcServer.Serve(lis)
	h.CheckError(errorMsg, "server listen register server service")
}


func (s *Server) IncrementLamportTimestamp(clientTimeStamp int64) {
	s.lamportTimeStamp = (h.Max(clientTimeStamp, s.lamportTimeStamp) + 1)
}
