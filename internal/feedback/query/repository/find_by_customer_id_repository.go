package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"ufriend-cx-dashboard-server/internal/feedback/model"
)

func (r *feedbackQueryRepository) FindByCustomerID(ctx context.Context, customerId string) ([]*model.Feedback, error) {
	objID, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return nil, err
	}

	cursor, err := r.collection.Find(ctx, bson.M{"customer_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var feedbacks []*model.Feedback
	if err := cursor.All(ctx, &feedbacks); err != nil {
		return nil, err
	}

	return feedbacks, nil
}
