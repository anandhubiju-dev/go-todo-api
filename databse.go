package main

import (
	"context"
	"todo-api/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func connectDB(cfg config.Config) (*mongo.Client, *mongo.Collection, error) {
	ctx := context.Background()

	mongoURI := cfg.MongoURI

	client, err := mongo.Connect(
		options.Client().ApplyURI(mongoURI),
	)

	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	db := client.Database("todo_db")
	todoCollection := db.Collection("todos")

	return client, todoCollection, nil
}
