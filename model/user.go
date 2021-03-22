package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//User struct declaration
type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	EncodedID string             `bson:"encoded_id" json:"encoded_id"`
}
