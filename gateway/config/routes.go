package config

import (
	pb "gateway/proto"
	"gateway/service"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	userConnection, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	userGrpcClient := pb.NewUserServiceClient(userConnection)
	userController := service.Service{UserClient: userGrpcClient}

	e.POST("/users", userController.Register)
	e.POST("/users/login", userController.Login)
	e.GET("/users/:id", userController.GetUserById)
	e.PUT("/users/:id", userController.UpdateUser)

	return e
}
