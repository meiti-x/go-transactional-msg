package main

import (
	"github.com/meiti-x/go-transactional-msg/api"
	"log"
	"net"

	"github.com/meiti-x/go-transactional-msg/pb"
	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterFSServer(grpcServer, api.NewFsGRPCApi("./files"))

	log.Println("Starting gRPC Server on port 8080")

	log.Fatal(grpcServer.Serve(l))
}
