package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Collection {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@postit-mongo:27017")

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client.Database("postit-db").Collection("posts")
}
