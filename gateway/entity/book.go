package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Author        string             `bson:"author" json:"author"`
	PublishedDate string             `bson:"published_date" json:"published_date"`
	Status        string             `bson:"status" json:"status"`
}
