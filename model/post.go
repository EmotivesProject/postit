package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	User    string             `bson:"user json:"user"`
	Message string             `bson:"message json:"message"`
	Created time.Time          `bson:"created json:"created"`
}
