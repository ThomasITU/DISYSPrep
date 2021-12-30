package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	h "github.com/ThomasITU/DISYSPrep/HelperMethod"
	"github.com/ThomasITU/DISYSPrep/Proto"
	"google.golang.org/grpc"
)

const (
	SERVER_PORT        = 5001
	FRONT_END_LOG_FILE = "frontEndLog"
)

type FrontEnd struct {
	Proto.UnimplementedProtoServiceServer
	connectedUsers     []int
	replicaServerPorts map[int]bool
	arbiter            sync.Mutex
}

func main() {

	//init
	replicaServers := make(map[int]bool)
	users := make([]int, 0)

	frontEnd := FrontEnd{connectedUsers: users, replicaServerPorts: replicaServers, arbiter: sync.Mutex{}}
	go listen(&frontEnd)
	for {
		// begin searching for servers every 5 seconds
		frontEnd.FindActiveServers()
		time.Sleep(5 * time.Second)
	}
}

func (fe *FrontEnd) JoinService(ctx context.Context, request *Proto.JoinRequest) (*Proto.Response, error) {
	var msg string
	userId := int(request.GetUserId())
	fe.arbiter.Lock()

	// check if user id already exist in the array
	for _, user := range fe.connectedUsers {
		if user == userId {
			msg = fmt.Sprintf("A user with id: %v has already joined", userId)
			break
		}
	}

	//add userid to slice
	if msg == "" {
		fe.connectedUsers = append(fe.connectedUsers, userId)
		msg = fmt.Sprintf("Welcome user: %v", userId)
	}
	fe.arbiter.Unlock()
	return &Proto.Response{Msg: msg}, nil
}

// get value grpc method logic
func (fe *FrontEnd) GetValue(ctx context.Context, request *Proto.GetRequest) (*Proto.Value, error) {
	values := make(map[h.Value]int)
	for port, alive := range fe.replicaServerPorts {
		if alive {
			client, state := h.ConnectToPort(port)
			if state == "alive" {
				response, _ := client.GetValue(ctx, &Proto.GetRequest{})
				value := h.Value{Value: response.GetCurrentValue(), UserId: response.GetUserId()}
				temp := values[value]
				values[value] = (temp + 1)
			}
		}
	}
	currentReplicas := 0
	for _, votes := range values {
		currentReplicas += votes
	}
	for value, votes := range values {
		if votes > currentReplicas/2 {
			return &Proto.Value{CurrentValue: value.Value, UserId: value.UserId}, nil
		}
	}

	return nil, errors.New("replicas couldn't agree on one value")
}

// set value grpc method logic
func (fe *FrontEnd) SetValue(ctx context.Context, request *Proto.SetRequest) (*Proto.Response, error) {
	failedUpdates := 0
	var msg string
	for port, alive := range fe.replicaServerPorts {
		if alive {
			client, status := h.ConnectToPort(port)
			if status == "alive" {
				response, err := client.SetValue(ctx, request)
				if err != nil || !strings.Contains(response.Msg, "Updated") {
					failedUpdates++
				}
			}
		}
	}
	if failedUpdates < (1+len(fe.replicaServerPorts))/2 {
		msg = fmt.Sprintf("updated more than half of the replicas with the value %v by user: %v", request.GetRequestedValue(), request.GetUserId())
	}else {
		msg = "failed to update more then half of the replicas"
	}
	h.Logger(msg, FRONT_END_LOG_FILE)
	return &Proto.Response{Msg: msg}, nil
}

func (fe *FrontEnd) FindActiveServers() {
	for i := 0; i < h.MAX_REPLICAS; i++ {
		serverPort := SERVER_PORT + i
		_, status := h.ConnectToPort(serverPort)
		if status == "alive" {
			fmt.Printf("found alive server at port : %v\n", serverPort)
			fe.replicaServerPorts[serverPort] = true
		} else if status == "unknown" {
			fe.replicaServerPorts[serverPort] = false
		}
	}

}

func listen(fe *FrontEnd) {
	//listen on port
	lis, err := net.Listen("tcp", h.FRONT_END_ADDRESS)
	h.CheckError(err, "server setup net.listen")
	defer lis.Close()

	// register server this is a blocking call
	grpcServer := grpc.NewServer()
	Proto.RegisterProtoServiceServer(grpcServer, fe)
	errorMsg := grpcServer.Serve(lis)
	h.CheckError(errorMsg, "server listen register server service")
}
