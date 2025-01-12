package service

import (
	"context"
	"gateway/entity"
	pb "gateway/proto"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service struct {
	UserClient pb.UserServiceClient
}

func (s *Service) Register(c echo.Context) error {
	var payload entity.UserInput

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userPB := &pb.RegisterUserRequest{
		Username: payload.Username,
		Password: payload.Password,
	}

	_, err := s.UserClient.RegisterUser(c.Request().Context(), userPB)

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
}

func (s *Service) GetUserById(c echo.Context) error {
	userID := c.Param("id")
	token := c.Request().Header.Get("Authorization")

	userPB := &pb.GetUserByIdRequest{Id: userID}

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := s.UserClient.GetUserById(ctx, userPB)

	if err != nil {
		st, _ := status.FromError(err)

		errMessage := st.Message()

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Message: errMessage,
		})
	}

	return c.JSON(http.StatusOK, entity.ResponseOK{
		Message: "User found",
		Data:    res,
	})
}

func (s *Service) Login(c echo.Context) error {
	var payload entity.UserInput

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userPB := &pb.LoginUserRequest{
		Username: payload.Username,
		Password: payload.Password,
	}

	res, err := s.UserClient.LoginUser(c.Request().Context(), userPB)

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
		Message: "Login success",
		Data:    res,
	})
}

func (s *Service) UpdateUser(c echo.Context) error {
	var payload entity.UpdateUserPayload
	token := c.Request().Header.Get("Authorization")

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	userPB := &pb.UpdateUserRequest{
		Id:       c.Param("id"),
		Password: payload.Password,
	}

	md := metadata.Pairs("authorization", token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := s.UserClient.UpdateUser(ctx, userPB)

	if err != nil {
		st, _ := status.FromError(err)

		errMessage := st.Message()

		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Message: errMessage,
		})
	}

	return c.JSON(http.StatusOK, entity.ResponseOK{
		Message: "User updated",
		Data:    res,
	})
}
