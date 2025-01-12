package repository

import (
	"context"
	"fmt"
	"userservice/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	RegisterUser(payload entity.UserInput) (*entity.User, error)
	LoginUser(username, password string) (*entity.User, error)
	GetUserById(userID string) (*entity.User, error)
	UpdateUser(userID, password string) (*entity.User, error)
	DeleteUser(userID string) error
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

	if userExists := ur.collection.FindOne(context.Background(), bson.M{"username": payload.Username}).Err(); userExists == nil {
		return nil, fmt.Errorf("username %s already exists", payload.Username)
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

func (ur *userRepository) LoginUser(username, password string) (*entity.User, error) {
	var user entity.User

	err := ur.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

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
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}

		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) UpdateUser(userID, password string) (*entity.User, error) {
	var user entity.User

	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return nil, err
	}

	_, err = ur.GetUserById(userID)

	if err != nil {
		return nil, err
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"password": string(bcryptPassword),
		},
	}

	err = ur.collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, update).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) DeleteUser(userID string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		return err
	}

	_, err = ur.GetUserById(userID)

	if err != nil {
		return err
	}

	_, err = ur.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})

	if err != nil {
		return err
	}

	return nil
}
