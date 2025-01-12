package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserPayload struct {
	Password string `json:"password"`
}
