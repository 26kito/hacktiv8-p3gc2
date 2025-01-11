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
		var payload entity.UserInput

		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		userPB := &pb.RegisterUserRequest{
			Username: payload.Username,
			Password: payload.Password,
		}

		_, err := userGrpcClient.RegisterUser(c.Request().Context(), userPB)

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

	e.GET("/users/:id", func(c echo.Context) error {
		userID := c.Param("id")

		userPB := &pb.GetUserByIdRequest{Id: userID}

		res, err := userGrpcClient.GetUserById(c.Request().Context(), userPB)

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

		return c.JSON(http.StatusOK, entity.ResponseOK{
			Message: "User found",
			Data:    res,
		})
	})

	return e
}
