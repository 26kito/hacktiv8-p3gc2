package repository

import (
	"bookservice/entity"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	GetAllBook() ([]entity.Book, error)
	InsertBook(payload entity.InsertBookRequest) (*entity.Book, error)
	GetBookById(id string) (*entity.Book, error)
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

func (br *bookRepository) InsertBook(payload entity.InsertBookRequest) (*entity.Book, error) {
	newBook := entity.Book{
		ID:            primitive.NewObjectID(),
		Title:         payload.Title,
		Author:        payload.Author,
		PublishedDate: payload.PublishedDate,
		Status:        payload.Status,
	}

	_, err := br.collection.InsertOne(context.Background(), newBook)

	if err != nil {
		return nil, err
	}

	return &newBook, nil
}

func (br *bookRepository) GetBookById(id string) (*entity.Book, error) {
	var book entity.Book

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = br.collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&book)

	if err != nil {
		return nil, err
	}

	return &book, nil
}
