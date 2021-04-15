package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//User struct declaration
type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id"`
	Username string             `bson:"username" json:"username"`
}
