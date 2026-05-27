package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"ufriend-cx-dashboard-server/internal/feedback/model"
)

func (r *feedbackCommandRepository) Create(ctx context.Context, feedback *model.Feedback) error {
	// Look up customer's branch to keep denormalized branch field correct
	customersCol := r.collection.Database().Collection("customers")
	
	var customer struct {
		Branch string `bson:"branch"`
	}
	err := customersCol.FindOne(ctx, bson.M{"_id": feedback.CustomerId}).Decode(&customer)
	if err == nil {
		feedback.Branch = customer.Branch
	}

	_, err = r.collection.InsertOne(ctx, feedback)
	return err
}