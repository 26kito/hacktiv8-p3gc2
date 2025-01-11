package repository

import (
	"context"
	"userservice/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(payload entity.UserInput) (*entity.User, error)
	GetUserById(userID string) (*entity.User, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *userRepository {
	return &userRepository{collection}
}

func (ur *userRepository) RegisterUser(payload entity.UserInput) (*entity.User, error) {
	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.MinCost)

	if err != nil {
		return nil, err
	}

	user := entity.User{
		ID:       primitive.NewObjectID(),
		Username: payload.Username,
		Password: string(bcryptPassword),
	}

	_, err = ur.collection.InsertOne(context.Background(), user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserById(userID string) (*entity.User, error) {
	var user entity.User

	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return nil, err
	}

	err = ur.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
