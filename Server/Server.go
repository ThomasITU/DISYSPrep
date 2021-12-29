package main

import (
	"fmt"
	"net"
	"strconv"

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
	port int
}

func main() {

	//init
	server := Server{port: SERVER_PORT}

	listen(&server)
	fmt.Println("main has ended")
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
	h.CheckError(errorMsg, "server listen registerservice")
}
