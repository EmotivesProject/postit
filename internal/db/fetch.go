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
	result := findByID(hex, PostCollection, mongoDB)
	err = result.Decode(&post)

	return post, err
}

func FindCommentById(commentID string, mongoDB *mongo.Database) (model.Comment, error) {
	var comment model.Comment
	hex, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return comment, err
	}
	result := findByID(hex, CommentCollection, mongoDB)
	err = result.Decode(&comment)

	return comment, err
}

func FindLikeById(likeID string, mongoDB *mongo.Database) (model.Like, error) {
	var like model.Like
	hex, err := primitive.ObjectIDFromHex(likeID)
	if err != nil {
		return like, err
	}
	result := findByID(hex, LikeCollection, mongoDB)
	err = result.Decode(&like)

	return like, err
}

func findByID(id primitive.ObjectID, collection string, mongoDB *mongo.Database) *mongo.SingleResult {
	filter := bson.M{"_id": id}
	databaseCollection := mongoDB.Collection(collection)
	return databaseCollection.FindOne(context.TODO(), filter)
}

func FindByLikeIDS(UserLikes []primitive.ObjectID, mongoDB *mongo.Database) ([]model.Like, error) {
	var likes []model.Like
	query := bson.M{"_id": bson.M{"$in": UserLikes}}
	collection := mongoDB.Collection(LikeCollection)

	cur, err := collection.Find(context.TODO(), query)
	if err != nil {
		return likes, err
	}

	err = cur.All(context.TODO(), &likes)
	return likes, err
}
