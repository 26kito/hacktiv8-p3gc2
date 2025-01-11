package main

import (
	"log"
	"net"
	"userservice/pb"

	"google.golang.org/grpc"
)

type AuthService struct {
	pb.UnimplementedUserServiceServer
}

func main() {
	listen, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterUserServiceServer(server, &AuthService{})

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	log.Println("Starting gRPC server on port :50051")
}
