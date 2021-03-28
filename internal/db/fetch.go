package db

import (
	"context"
	"postit/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindUser(username string, mongoDB *mongo.Database) (model.User, error) {
	user := model.User{}
	filter := bson.D{primitive.E{Key: "username", Value: username}}

	usersCollection := mongoDB.Collection(UserCollection)
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	return user, err
}

func FindPostById(postID string, mongoDB *mongo.Database) (model.Post, error) {
	var post model.Post
	hex, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return post, err
	}
	filter := bson.M{"_id": hex}

	postsCollection := mongoDB.Collection(PostCollection)
	err = postsCollection.FindOne(context.TODO(), filter).Decode(&post)
	return post, err
}
