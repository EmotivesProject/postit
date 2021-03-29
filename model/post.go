package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID           primitive.ObjectID   `bson:"_id" json:"id"`
	User         primitive.ObjectID   `bson:"user" json:"user"`
	Message      string               `bson:"message" json:"message"`
	UserLikes    []primitive.ObjectID `bson:"user_likes" json:"user_likes"`
	UserComments []primitive.ObjectID `bson:"user_comments" json:"user_comments"`
	Created      time.Time            `bson:"created" json:"created"`
	Active       bool                 `bson:"active" json:"active"`
}
