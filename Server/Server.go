package main

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	h "github.com/ThomasITU/DISYSPrep/HelperMethod"

	"google.golang.org/grpc"

	"github.com/ThomasITU/DISYSPrep/Proto"
)

const (
	SERVER_PORT     = 5000
	SERVER_LOG_FILE = "serverLog"
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
	serverPort := FindFreePort()
	if serverPort == -1 { // if no free port -1
		fmt.Printf("Can't start more than %v", h.MAX_REPLICAS)
		return
	}
	server := Server{port: serverPort, latestValue: initValue, arbiter: sync.Mutex{}}
	fmt.Printf("Succesfully got port: %v", server.port) // sanity checks

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

func (s *Server) JoinService(ctx context.Context, request *Proto.JoinRequest) (*Proto.Response, error) {
	userId := request.GetUserId()
	if userId == -1 {
		return &Proto.Response{Msg: "alive"}, nil
	} else {
		msg := fmt.Sprintf("Welcome to our marvelous service user: %v ", userId)
		return &Proto.Response{Msg: msg}, nil
	}
}

// connect to ports until a free port is found
func FindFreePort() int {
	for i := 1; i < (h.MAX_REPLICAS + 1); i++ {
		serverPort := SERVER_PORT + i
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(serverPort), grpc.WithTimeout(time.Millisecond*250), grpc.WithInsecure())
		if err == nil {
			defer conn.Close()
			ctx := context.Background()
			client := Proto.NewProtoServiceClient(conn)
			response, _ := client.JoinService(ctx, &Proto.JoinRequest{UserId: -1})
			if response.GetMsg() == "alive" {
				continue
			} else {
				return serverPort
			}
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
