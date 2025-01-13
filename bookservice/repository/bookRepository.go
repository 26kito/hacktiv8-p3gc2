package repository

import (
	"bookservice/entity"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	GetAllBook() ([]entity.Book, error)
	InsertBook(payload entity.InsertBookRequest) (*entity.Book, error)
	GetBookById(id string) (*entity.Book, error)
	UpdateBook(payload entity.UpdateBookRequest) (*entity.Book, error)
	DeleteBook(id string) error
	BorrowBook(payload entity.BorrowBookRequest) (*entity.BorrowBook, error)
	ReturnBook(payload entity.ReturnBookRequest) error
}

type bookRepository struct {
	// collection *mongo.Collection
	booksCollection         *mongo.Collection
	borrowedBooksCollection *mongo.Collection
}

func NewBookRepository(booksCollection, borrowedBooksCollection *mongo.Collection) *bookRepository {
	return &bookRepository{
		booksCollection:         booksCollection,
		borrowedBooksCollection: borrowedBooksCollection,
	}
}

func (br *bookRepository) GetAllBook() ([]entity.Book, error) {
	var books []entity.Book

	cursor, err := br.booksCollection.Find(context.Background(), bson.M{})

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

	_, err := br.booksCollection.InsertOne(context.Background(), newBook)

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

	err = br.booksCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&book)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (br *bookRepository) UpdateBook(payload entity.UpdateBookRequest) (*entity.Book, error) {
	objID, err := primitive.ObjectIDFromHex(payload.ID)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"title":          payload.Title,
			"author":         payload.Author,
			"published_date": payload.PublishedDate,
			"status":         payload.Status,
		},
	}

	_, err = br.booksCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		return nil, err
	}

	return &entity.Book{
		ID:            objID,
		Title:         payload.Title,
		Author:        payload.Author,
		PublishedDate: payload.PublishedDate,
		Status:        payload.Status,
	}, nil
}

func (br *bookRepository) DeleteBook(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = br.booksCollection.DeleteOne(context.Background(), bson.M{"_id": objID})

	if err != nil {
		return err
	}

	return nil
}

func (br *bookRepository) BorrowBook(payload entity.BorrowBookRequest) (*entity.BorrowBook, error) {
	bookID, err := primitive.ObjectIDFromHex(payload.BookID)
	if err != nil {
		return nil, err
	}

	userID, err := primitive.ObjectIDFromHex(payload.UserID)
	if err != nil {
		return nil, err
	}

	book, err := br.GetBookById(payload.BookID)
	if err != nil {
		return nil, err
	}

	if book.Status != "available" {
		return nil, fmt.Errorf("book is not available")
	}

	newBorrowBook := entity.BorrowBook{
		ID:         primitive.NewObjectID(),
		BookID:     bookID,
		UserID:     userID,
		BorrowDate: payload.BorrowDate,
		ReturnDate: "",
	}

	_, err = br.borrowedBooksCollection.InsertOne(context.Background(), newBorrowBook)
	if err != nil {
		return nil, err
	}

	_, err = br.booksCollection.UpdateOne(context.Background(), bson.M{"_id": bookID}, bson.M{"$set": bson.M{"status": "borrowed"}})

	return &newBorrowBook, nil
}

func (br *bookRepository) BorrowBookDetail(id string) (*entity.BorrowBook, error) {
	var borrowBook entity.BorrowBook

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = br.borrowedBooksCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&borrowBook)

	if err != nil {
		return nil, err
	}

	return &borrowBook, nil
}

func (br *bookRepository) ReturnBook(payload entity.ReturnBookRequest) error {
	objID, err := primitive.ObjectIDFromHex(payload.ID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"return_date": payload.ReturnDate,
		},
	}

	_, err = br.borrowedBooksCollection.UpdateOne(context.Background(), bson.M{"_id": objID}, update)

	if err != nil {
		return err
	}

	borrowedBook, err := br.BorrowBookDetail(payload.ID)
	if err != nil {
		return err
	}

	_, err = br.booksCollection.UpdateOne(context.Background(), bson.M{"_id": borrowedBook.BookID}, bson.M{"$set": bson.M{"status": "available"}})

	return nil
}
