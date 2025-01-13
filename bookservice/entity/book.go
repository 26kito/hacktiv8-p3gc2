package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Title         string             `bson:"title"`
	Author        string             `bson:"author"`
	PublishedDate string             `bson:"published_date"`
	Status        string             `bson:"status"`
}

type InsertBookRequest struct {
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	Status        string `json:"status"`
}

type UpdateBookRequest struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"published_date"`
	Status        string `json:"status"`
}
