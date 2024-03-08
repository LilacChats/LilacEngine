package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeConnection(URL string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(URL)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	} else {
		return client, nil
	}
}
