package repository

import (
	"bookservice/entity"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	GetAllBook() ([]entity.Book, error)
}

type bookRepository struct {
	collection *mongo.Collection
}

func NewBookRepository(collection *mongo.Collection) *bookRepository {
	return &bookRepository{collection}
}

func (br *bookRepository) GetAllBook() ([]entity.Book, error) {
	var books []entity.Book

	cursor, err := br.collection.Find(context.Background(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var book entity.Book
		cursor.Decode(&book)
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
