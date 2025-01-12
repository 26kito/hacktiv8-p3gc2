package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"userservice/config"
	pb "userservice/proto"
	"userservice/repository"
	"userservice/service"

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
		"/user.UserService/GetUserById": true,
		"/user.UserService/UpdateUser":  true,
	}

	// Check if the route requires authentication
	if protectedRoutes[info.FullMethod] {
		// Extract metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md["authorization"][0]
		if authHeader == "" {
			return nil, status.Errorf(codes.Unauthenticated, "missing token")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token format")
		}

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}

			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			return nil, fmt.Errorf("invalid or expired token")
		}

		// Store the claims in the context for further use
		claims := token.Claims.(jwt.MapClaims)
		const userClaimsKey contextKey = "user"
		ctx = context.WithValue(ctx, userClaimsKey, claims)

		// Call the handler
		return handler(ctx, req)
	}

	// Call the handler
	return handler(ctx, req)
}

func main() {
	listen, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	godotenv.Load()

	db, err := config.Connect(context.Background())

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a gRPC server with the interceptor
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor),
	)

	userRepository := repository.NewUserRepository(db)
	userservice := service.NewUserService(userRepository)

	pb.RegisterUserServiceServer(server, userservice)

	log.Println("Starting gRPC server on port :50051")

	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
