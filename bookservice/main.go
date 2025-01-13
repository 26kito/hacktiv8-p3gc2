package main

import (
	"bookservice/config"
	"bookservice/repository"
	"bookservice/service"
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	pb "bookservice/proto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

var JWTSecret = []byte("hacktiv8p3gc2")

func UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Define protected routes
	protectedRoutes := map[string]bool{
		"/bookservice.BookService/BorrowBook": true,
		"/bookservice.BookService/ReturnBook": true,
	}

	// Check if the route requires authentication
	if !protectedRoutes[info.FullMethod] {
		return handler(ctx, req)
	}

	// Extract metadata
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok || len(md["authorization"]) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing authorization metadata")
	}

	// Extract the token
	authHeader := md["authorization"][0]
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid or expired token")
	}

	// Store claims in context
	claims := token.Claims.(jwt.MapClaims)
	ctx = context.WithValue(ctx, "user", claims)

	// Proceed with the handler
	return handler(ctx, req)
}

func main() {
	listen, err := net.Listen("tcp", ":50052")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	godotenv.Load()

	booksCollection, borrowedBooksCollection, err := config.Connect(context.Background())

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a gRPC server with the interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor),
	)

	// Register the service with the server
	bookRepository := repository.NewBookRepository(booksCollection, borrowedBooksCollection)
	bookService := service.NewBookService(bookRepository)

	pb.RegisterBookServiceServer(server, bookService)
	log.Println("Starting gRPC server on port :50052")

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
