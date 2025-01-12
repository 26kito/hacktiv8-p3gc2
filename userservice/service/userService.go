package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"userservice/entity"
	pb "userservice/proto"
	"userservice/repository"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/status"
)

type UserService struct {
	UserRepository repository.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

func (us *UserService) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	payload := entity.UserInput{
		Username: req.Username,
		Password: req.Password,
	}

	if err := validateRegisterPayload(payload); err != nil {
		return nil, status.Error(400, err.Error())
	}

	res, err := us.UserRepository.RegisterUser(payload)

	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		Id: res.ID.Hex(),
	}, nil
}

func (us *UserService) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.UserResponse, error) {
	username := req.Username
	password := req.Password

	res, err := us.UserRepository.LoginUser(username, password)

	if err != nil {
		return nil, err
	}

	tokenString, err := generateJWTToken(res)

	if err != nil {
		return nil, err
	}

	return &pb.UserResponse{
		Id:       res.ID.Hex(),
		Username: res.Username,
		Token:    tokenString,
	}, nil
}

func (us *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	userID := req.Id

	res, err := us.UserRepository.GetUserById(userID)

	if err != nil {
		return nil, err
	}

	return &pb.GetUserByIdResponse{
		Id:       res.ID.Hex(),
		Username: res.Username,
	}, nil
}

func validateRegisterPayload(payload entity.UserInput) error {
	if payload.Username == "" {
		return fmt.Errorf("username is required")
	}

	if strings.Contains(payload.Username, " ") {
		return fmt.Errorf("username cannot contain spaces")
	}

	if len(payload.Username) < 5 || len(payload.Username) > 15 {
		return fmt.Errorf("username must be between 5 and 15 characters")
	}

	if payload.Password == "" {
		return fmt.Errorf("password is required")
	}

	if len(payload.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	return nil
}

func generateJWTToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte("hacktiv8p3gc2"))

	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}
