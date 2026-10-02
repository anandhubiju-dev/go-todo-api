package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Todo struct {
	ID          bson.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string        `json:"title" bson:"title" binding:"required"`
	Description string        `json:"description" bson:"description"`
	Completed   bool          `json:"completed" bson:"completed"`
}
