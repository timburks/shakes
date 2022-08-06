package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/timburks/shakes/rpc"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type queryServer struct {
	rpc.UnimplementedQueryServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	rpc.RegisterQueryServer(grpcServer, &queryServer{})
	grpcServer.Serve(lis)
}
