package db

import (
	"context"
	"log"
	"time"

	"github.com/TomBowyerResearchProject/common/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Database
)

func Connect() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin@mongo:27017")

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Connected to MongoDB!")

	db = client.Database("postit-db")
}

func GetDatabase() *mongo.Database {
	return db
}
