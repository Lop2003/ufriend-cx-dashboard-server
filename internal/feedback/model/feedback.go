package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feedback struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	CustomerId primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	Rating     int                `json:"rating" bson:"rating"`
	Comment    string             `json:"comment" bson:"comment"`
	Category   string             `json:"category" bson:"category"`
	Sentiment  string             `json:"sentiment" bson:"sentiment"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}