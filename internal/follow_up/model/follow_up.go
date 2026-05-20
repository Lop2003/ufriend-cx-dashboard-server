package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FollowUp struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	CustomerId primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	Type       string             `json:"type" bson:"type"`
	Note       string             `json:"note" bson:"note"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}