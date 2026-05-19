package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Phone      string             `json:"phone" bson:"phone"`
	Product    string             `json:"product" bson:"product"`
	Branch     string             `json:"branch" bson:"branch"`
	PlanMonths int                `json:"plan_months" bson:"plan_months"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}