package db_url

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	db     *mongo.Database
}

// NewRepository is a constructor for the Repository
func NewRepository(connectionString, dbName string) (*Repository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &Repository{client: client, db: db}, nil
}

// Find is a method on the Repository struct.
func (r *Repository) Find(shortURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.db.Collection("urls")

	filter := bson.D{{Key: "shortUrl", Value: shortURL}}

	options := options.FindOne().SetProjection(bson.M{"url": 1, "_id": 0})

	var result bson.M
	err := collection.FindOne(ctx, filter, options).Decode(&result)
	if err != nil {
		return "", err
	}

	url, ok := result["url"].(string)
	if !ok {
		return "", errors.New("unable to cast url to string")
	}

	baseURL := os.Getenv("BASE_URL")

	return baseURL + url, nil
}

// Save is a method on the Repository struct.
func (r *Repository) Save(shortURL string, longURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.db.Collection("urls")

	doc := URLDocument{ShortURL: shortURL, URL: longURL}

	_, err := collection.InsertOne(ctx, doc)
	return err
}

// Disconnect is a method to disconnect from the database.
func (r *Repository) Disconnect() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}
}
