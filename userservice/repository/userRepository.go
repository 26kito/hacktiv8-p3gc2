package repository

import (
	"context"
	"userservice/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{collection}
}

func (r *UserRepository) CreateUser(username, password string) error {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		return err
	}

	user := entity.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: string(bcryptPassword),
	}

	_, err = r.collection.InsertOne(context.Background(), user)

	if err != nil {
		return err
	}

	return nil
}
