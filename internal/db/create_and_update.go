package db

import (
	"context"
	"encoding/json"
	"io"
	"postit/model"
	"postit/pkg/postit_messages"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func CreatePost(body io.ReadCloser, userID primitive.ObjectID, mongoDB *mongo.Database) (*mongo.InsertOneResult, error) {
	post := &model.Post{}
	err := json.NewDecoder(body).Decode(post)
	if err != nil {
		return &mongo.InsertOneResult{}, postit_messages.ErrFailedDecoding
	}

	post.Created = time.Now()
	post.User = userID
	post.ID = primitive.NewObjectID()
	post.Active = true

	return insetIntoCollection(PostCollection, post, mongoDB)
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

func CreateUser(body io.ReadCloser, mongoDB *mongo.Database) (*model.User, error) {
	user := &model.User{}
	err := json.NewDecoder(body).Decode(user)
	if err != nil {
		return user, err
	}
	user.ID = primitive.NewObjectID()

	_, err = insetIntoCollection(UserCollection, user, mongoDB)
	return user, err
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
