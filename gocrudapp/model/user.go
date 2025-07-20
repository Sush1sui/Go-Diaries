package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID      bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string        `bson:"name,omitempty" json:"name,omitempty"`
	Age     int           `bson:"age,omitempty" json:"age,omitempty"`
	Country string        `bson:"country,omitempty" json:"country,omitempty"`
}
