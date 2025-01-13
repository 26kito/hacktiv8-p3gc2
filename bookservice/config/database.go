package config

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect(ctx context.Context) (*mongo.Collection, *mongo.Collection, error) {
	MONGO_URI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))

	if err != nil {
		return nil, nil, err
	}

	booksCollection := client.Database("hacktiv8-p3gc2").Collection("books")
	borrowedBooksCollection := client.Database("hacktiv8-p3gc2").Collection("borrowed_books")

	return booksCollection, borrowedBooksCollection, nil
}
