package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Tasks struct{
	Id 			primitive.ObjectID		`json:"id" bson:"_id,omitempty"`
	Title 		string					`json:"title" bson:"title"`
	Description	string					`json:"description" bson:"description"`
	Status		string					`json:"status" bson:"status"`
	CreatedAt	time.Time				`json:"created_at" bson:"created_at"`
	UpdatedAt 	time.Time				`json:"updated-at" bson:"updated_at"`
}