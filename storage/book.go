package storage

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"isbn-locator/library"
)

type BookStorage struct {
	collection *mongo.Collection
}

func NewBookStorage(database *mongo.Database) *BookStorage {
	collection := database.Collection("books")
	return &BookStorage{
		collection: collection,
	}
}

var _ library.BookRepository = (*BookStorage)(nil)

func (storage *BookStorage) Fetch(ctx context.Context, isbn string) (library.Book, error) {
	log.Printf("Document requested with ISBN: %v\n", isbn)
	var book library.Book
	filter := bson.D{{Key: "isbn", Value: isbn}}
	err := storage.collection.FindOne(ctx, filter).Decode(&book)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return book, library.ErrBookNotFound
		}
		return book, err
	}
	return book, nil
}

func (storage *BookStorage) Store(ctx context.Context, b library.Book) error {
	insertResult, err := storage.collection.InsertOne(ctx, b)
	if err != nil {
		return err
	}
	log.Printf("Document %s inserted with ISBN: %v\n", insertResult.InsertedID, b.ISBN)
	return nil
}

func (storage *BookStorage) Update(ctx context.Context, isbn string, b library.Book) error {
	count := int64(0)
	if isbn == b.ISBN {
		filter := bson.D{{Key: "isbn", Value: isbn}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: b.Title}, {Key: "author", Value: b.Author}, {Key: "year", Value: b.Year}}}}
		updateResult, err := storage.collection.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
		count = updateResult.ModifiedCount
	}
	log.Printf("Documents updated: %v\n", count)
	return nil
}

func (storage *BookStorage) Remove(ctx context.Context, isbn string) error {
	filter := bson.D{{Key: "isbn", Value: isbn}}
	delResult, err := storage.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	log.Printf("Documents deleted: %v\n", delResult.DeletedCount)
	return nil
}
