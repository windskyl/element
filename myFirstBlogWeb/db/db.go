package db

import (
	"context"
	"myFirstBlogWeb/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(config.MongoDBURI)
	var err error
	Client, err = mongo.Connect(ctx, clientOpts)
	return err
}

func GetCollection(name string) *mongo.Collection {
	return Client.Database(config.DatabaseName).Collection(name)
}

func CloseDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	Client.Disconnect(ctx)
}
