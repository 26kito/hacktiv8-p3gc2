package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	Author        string             `bson:"author"`
	PublishedDate string             `bson:"published_date"`
	Status        string             `bson:"status"`
}
