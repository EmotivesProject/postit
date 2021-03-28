package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetRawResponseFromAggregate(collectionName string, pipeline mongo.Pipeline, mongoDB *mongo.Database) ([]bson.M, error) {
	var rawCollection []bson.M
	collection := mongoDB.Collection(collectionName)
	cur, err := collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return rawCollection, err
	}

	err = cur.All(context.TODO(), &rawCollection)
	return rawCollection, err
}
