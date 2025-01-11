package service

import (
	"userservice/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository}
}

func (s *UserService) CreateUser(username, password string) error {
	err := s.UserRepository.CreateUser(username, password)

	if err != nil {
		return err
	}

	return nil
}
