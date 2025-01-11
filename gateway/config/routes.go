package config

import (
	"net/http"

	"gateway/entity"
	pb "gateway/proto"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func NewRouter() *echo.Echo {
	e := echo.New()

	userConnection, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		e.Logger.Fatalf("did not connect: %v", err)
	}

	userGrpcClient := pb.NewUserServiceClient(userConnection)

	e.POST("/users", func(c echo.Context) error {
		var user entity.UserInput

		if err := c.Bind(&user); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		userPB := &pb.CreateUserRequest{
			Username: user.Username,
			Password: user.Password,
		}

		_, err := userGrpcClient.CreateUser(c.Request().Context(), userPB)

		if err != nil {
			st, ok := status.FromError(err)

			errMessage := st.Message()

			if !ok {
				return c.JSON(http.StatusBadRequest, entity.ResponseError{
					Message: errMessage,
				})
			}

			return c.JSON(http.StatusBadRequest, entity.ResponseError{
				Message: errMessage,
			})
		}

		return c.JSON(http.StatusCreated, entity.ResponseOK{
			Message: "User created successfully",
			Data:    nil,
		})
	})

	return e
}
