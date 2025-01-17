package config

import (
	_ "gateway/docs"
	bookPb "gateway/proto/book"
	userPb "gateway/proto/user"
	"gateway/service"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	userConnection, err := grpc.NewClient("userservice:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	bookConnection, err := grpc.NewClient("bookservice:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	userGrpcClient := userPb.NewUserServiceClient(userConnection)
	userController := service.Service{UserClient: userGrpcClient}
	bookGrpcClient := bookPb.NewBookServiceClient(bookConnection)
	bookController := service.BookService{BookClient: bookGrpcClient}

	e.POST("/users", userController.Register)
	e.POST("/users/login", userController.Login)
	e.GET("/users/:id", userController.GetUserById)
	e.PUT("/users/:id", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)

	e.GET("/books", bookController.GetAllBook)
	e.POST("/books", bookController.InsertBook)
	e.GET("/books/:id", bookController.GetBookById)
	e.PUT("/books/:id", bookController.UpdateBook)
	e.DELETE("/books/:id", bookController.DeleteBook)

	e.POST("/books/:id/borrow", bookController.BorrowBook)
	e.POST("/books/:id/return", bookController.ReturnBook)

	e.GET("/cron/book-update-status", bookController.UpdateBookStatus)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}
