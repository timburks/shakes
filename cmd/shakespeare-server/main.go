package main

import (
	"log"
	"net"
	"os"

	"github.com/timburks/shakespeare/rpc"
	"google.golang.org/grpc"
)

type queryServer struct {
	rpc.UnimplementedQueryServer
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listening on port %s", port)
	grpcServer := grpc.NewServer()
	rpc.RegisterQueryServer(grpcServer, &queryServer{})
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
