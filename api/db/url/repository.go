package db_url

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/paoloposso/url_shrt/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client  *mongo.Client
	db      *mongo.Database
	timeout time.Duration
}

func NewRepository(connectionString, dbName string, c util.ConfigService) (*Repository, error) {

	timeout := c.GetMongoDbTimeOut()

	ctx, _ := context.WithTimeout(context.Background(), timeout)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	collection := db.Collection("urls")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, err
	}

	return &Repository{client: client, db: db}, nil
}

func (r *Repository) Find(shortURL string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), r.timeout)

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

	return url, nil
}

func (r *Repository) Save(shortURL string, longURL string) error {
	ctx, _ := context.WithTimeout(context.Background(), r.timeout)

	collection := r.db.Collection("urls")

	filter := bson.D{{Key: "url", Value: longURL}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "shortUrl", Value: shortURL},
		{Key: "url", Value: longURL},
	}}}

	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *Repository) Disconnect() {
	if err := r.client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}
