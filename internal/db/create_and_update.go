package db

import (
	"context"
	"encoding/json"
	"io"
	"postit/internal/postit_messages"
	"postit/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreatePost(body io.ReadCloser, userID primitive.ObjectID, mongoDB *mongo.Database) (*mongo.InsertOneResult, *model.Post, error) {
	post := &model.Post{}
	err := json.NewDecoder(body).Decode(post)
	if err != nil {
		return &mongo.InsertOneResult{}, post, postit_messages.ErrFailedDecoding
	}

	post.Created = time.Now()
	post.User = userID
	post.ID = primitive.NewObjectID()
	post.Active = true

	if !post.Validate() {
		return nil, post, postit_messages.ErrInvalid
	}

	mongoResult, err := insetIntoCollection(PostCollection, post, mongoDB)
	return mongoResult, post, err
}

func CreateComment(body io.ReadCloser, userID primitive.ObjectID, mongoDB *mongo.Database) (*model.Comment, error) {
	comment := &model.Comment{}
	err := json.NewDecoder(body).Decode(comment)
	if err != nil {
		return comment, postit_messages.ErrFailedDecoding
	}

	comment.ID = primitive.NewObjectID()
	comment.User = userID
	comment.Created = time.Now()
	comment.Active = true

	_, err = insetIntoCollection(CommentCollection, comment, mongoDB)
	return comment, err
}

func CreateUser(username string, mongoDB *mongo.Database) (*model.User, error) {
	user := model.User{
		Username: username,
	}
	user.ID = primitive.NewObjectID()

	_, err := insetIntoCollection(UserCollection, user, mongoDB)
	return &user, err
}

func CreateLike(userID primitive.ObjectID, mongoDB *mongo.Database) (*model.Like, error) {
	like := &model.Like{
		ID:      primitive.NewObjectID(),
		User:    userID,
		Created: time.Now(),
		Active:  true,
	}

	_, err := insetIntoCollection(LikeCollection, like, mongoDB)
	return like, err
}

func insetIntoCollection(collectionName string, document interface{}, mongoDB *mongo.Database) (*mongo.InsertOneResult, error) {
	collection := mongoDB.Collection(collectionName)
	return collection.InsertOne(context.TODO(), document)
}

func UpdatePost(post *model.Post, mongoDB *mongo.Database) error {
	postCollection := mongoDB.Collection(PostCollection)
	_, err := postCollection.ReplaceOne(context.TODO(), bson.M{"_id": post.ID}, post)
	return err
}

func UpdateComment(comment *model.Comment, mongoDB *mongo.Database) error {
	commentCollection := mongoDB.Collection(CommentCollection)
	_, err := commentCollection.ReplaceOne(context.TODO(), bson.M{"_id": comment.ID}, comment)
	return err
}

func UpdateLike(like *model.Like, mongoDB *mongo.Database) error {
	likeCollection := mongoDB.Collection(LikeCollection)
	_, err := likeCollection.ReplaceOne(context.TODO(), bson.M{"_id": like.ID}, like)
	return err
}
