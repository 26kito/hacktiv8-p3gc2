package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title         string             `bson:"title" json:"title"`
	Author        string             `bson:"author" json:"author"`
	PublishedDate string             `bson:"published_date" json:"published_date"`
	Status        string             `bson:"status" json:"status"`
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

type BorrowBook struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	BookID     primitive.ObjectID `bson:"book_id"`
	UserID     primitive.ObjectID `bson:"user_id"`
	BorrowDate string             `bson:"borrow_date"`
	ReturnDate string             `bson:"return_date"`
}

type BorrowBookRequest struct {
	BookID     string `json:"book_id"`
	UserID     string `json:"user_id"`
	BorrowDate string `json:"borrow_date"`
}

type ReturnBookRequest struct {
	ID         string `json:"id"`
	ReturnDate string `json:"return_date"`
}
