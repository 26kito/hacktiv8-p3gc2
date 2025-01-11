package service

import (
	"context"
	"fmt"
	"userservice/entity"
	pb "userservice/proto"
	"userservice/repository"

	"google.golang.org/grpc/status"
)

type UserService struct {
	UserRepository repository.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

func (us *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	payload := entity.UserInput{
		Username: req.Username,
		Password: req.Password,
	}

	if err := validateRegisterPayload(payload); err != nil {
		return nil, status.Error(400, err.Error())
	}

	res, err := us.UserRepository.CreateUser(payload)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id: res.ID.Hex(),
	}, nil
}

func validateRegisterPayload(payload entity.UserInput) error {
	if payload.Username == "" {
		return fmt.Errorf("username is required")
	}

	if payload.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}
