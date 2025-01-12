package main

import (
	"bookservice/config"
	"bookservice/repository"
	"bookservice/service"
	"context"
	"log"
	"net"

	pb "bookservice/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":50052")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	godotenv.Load()

	db, err := config.Connect(context.Background())

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a gRPC server with the interceptor
	server := grpc.NewServer()

	// Register the service with the server
	bookRepository := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepository)

	pb.RegisterBookServiceServer(server, bookService)

	log.Println("Starting gRPC server on port :50052")

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
