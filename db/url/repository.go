package db_url

import (
	"context"
	"errors"
	"log"
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

	collection := db.Collection("urls")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},  // Index key
		Options: options.Index().SetUnique(true), // Make it a unique index
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, err
	}

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

	// baseURL := os.Getenv("BASE_URL")

	return url, nil
}

// Save is a method on the Repository struct.
func (r *Repository) Save(shortURL string, longURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := r.db.Collection("urls")

	filter := bson.D{{Key: "url", Value: longURL}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "shortUrl", Value: shortURL},
		{Key: "url", Value: longURL},
	}}}

	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
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
