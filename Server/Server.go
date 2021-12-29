package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	h "../HelperMethod"
	"google.golang.org/grpc"

	"github.com/ThomasITU/DISYSPrep/Proto"
)

const (
	SERVER_PORT     = 5000
	SERVER_LOG_FILE = "serverLog"
)

type Server struct {
	Proto.UnimplementedProtoServiceServer
	port           int
	latestValue    Value
	connectedUsers []int
	arbiter        sync.Mutex
}

type Value struct {
	value  int64
	userId int64
}

func main() {

	//init
	initValue := Value{value: -1, userId: -1}
	users := make([]int, 0)
	lock := sync.Mutex{}
	server := Server{port: SERVER_PORT, latestValue: initValue, connectedUsers: users, arbiter: lock}

	listen(&server)
	fmt.Println("main has ended")
}

// JoinService grpc method logic
func (s *Server) JoinService(ctx context.Context, request *Proto.JoinRequest) (*Proto.Response, error) {
	s.arbiter.Lock()
	var msg string
	userId := int(request.GetUserId())

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
	s.arbiter.Unlock()
	return &Proto.Response{Msg: msg}, nil
}

// getvalue grpc method logic
func (s *Server) GetValue(ctx context.Context, request *Proto.GetRequest) (*Proto.Value, error) {
	value := Proto.Value{CurrentValue: s.latestValue.value, UserId: s.latestValue.userId}
	return &value, nil
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
