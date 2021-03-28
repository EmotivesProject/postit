package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	User    primitive.ObjectID `bson:"user" json:"user"`
	Message string             `bson:"message" json:"message"`
	Created time.Time          `bson:"created" json:"created"`
	Active  bool               `bson:"active" json:"active"`
}
