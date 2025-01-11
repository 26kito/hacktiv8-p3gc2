package repository

import (
	"context"
	"userservice/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(payload entity.UserInput) (*entity.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *userRepository {
	return &userRepository{collection}
}

func (r *userRepository) CreateUser(payload entity.UserInput) (*entity.User, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.MinCost)

	if err != nil {
		return nil, err
	}

	user := entity.User{
		ID:       primitive.NewObjectID(),
		Username: payload.Username,
		Password: string(bcryptPassword),
	}

	_, err = r.collection.InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
