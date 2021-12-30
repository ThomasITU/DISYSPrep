package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	h "github.com/ThomasITU/DISYSPrep/HelperMethod"

	"google.golang.org/grpc"

	"github.com/ThomasITU/DISYSPrep/Proto"
)

const (
	SERVER_PORT     = 5000
	SERVER_LOG_FILE = "serverLog"
	MAX_REPLICAS    = 5
)

type Server struct {
	Proto.UnimplementedProtoServiceServer
	port        int
	latestValue h.Value
	arbiter     sync.Mutex
}

func main() {

	//init
	initValue := h.Value{Value: -1, UserId: -1}
	lock := sync.Mutex{}
	serverPort := FindFreePort()
	if serverPort == -1 {
		fmt.Printf("Can't start more than %v", MAX_REPLICAS)
		return
	}

	server := Server{port: serverPort, latestValue: initValue, arbiter: lock}

	listen(&server)
	fmt.Println("main has ended")
}

// get value grpc method logic
func (s *Server) GetValue(ctx context.Context, request *Proto.GetRequest) (*Proto.Value, error) {
	value := Proto.Value{CurrentValue: s.latestValue.Value, UserId: s.latestValue.UserId}
	return &value, nil
}

// set value grpc method logic
func (s *Server) SetValue(ctx context.Context, request *Proto.SetRequest) (*Proto.Response, error) {
	s.arbiter.Lock()
	temp := s.latestValue
	s.latestValue = h.Value{Value: request.GetRequestedValue(), UserId: request.GetUserId()}
	msg := fmt.Sprintf("Updated the value: %v by %v to %v by %v ", temp.Value, temp.UserId, s.latestValue.Value, s.latestValue.UserId)
	h.Logger(msg, SERVER_LOG_FILE)
	s.arbiter.Unlock()
	return &Proto.Response{Msg: msg}, nil
}

// connect to ports until a free port is found
func FindFreePort() int {
	for i := 1; i < (MAX_REPLICAS + 1); i++ {
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(SERVER_PORT+i), grpc.WithInsecure(), grpc.WithBlock())
		if err == nil {
			defer conn.Close()
			return i + SERVER_PORT
		}
	}
	return -1
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
