package main

import (
	"context"
	"log"
	"net"
	"userservice/config"
	pb "userservice/proto"
	"userservice/repository"
	"userservice/service"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	godotenv.Load()

	server := grpc.NewServer()

	db, err := config.Connect(context.Background())

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepository := repository.NewUserRepository(db)
	userservice := service.NewUserService(userRepository)

	pb.RegisterUserServiceServer(server, userservice)

	log.Println("Starting gRPC server on port :50051")

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
