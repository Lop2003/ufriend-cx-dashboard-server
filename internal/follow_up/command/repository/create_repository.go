package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/follow_up/model"
)

func (r *followUpCommandRepository) Create(ctx context.Context, followUp *model.FollowUp) error {
	// Verify customer existence to maintain referential integrity
	customersCol := r.collection.Database().Collection("customers")
	count, err := customersCol.CountDocuments(ctx, bson.M{"_id": followUp.CustomerId})
	if err != nil {
		return err
	}
	if count == 0 {
		return mongo.ErrNoDocuments
	}

	_, err = r.collection.InsertOne(ctx, followUp)
	return err
}