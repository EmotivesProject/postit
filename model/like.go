package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Like struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Post    primitive.ObjectID `bson:"post" json:"post"`
	User    primitive.ObjectID `bson:"user" json:"user"`
	Created time.Time          `bson:"created" json:"created"`
	Active  bool               `bson:"active" json:"active"`
}
