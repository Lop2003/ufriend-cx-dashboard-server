package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"ufriend-cx-dashboard-server/internal/feedback/model"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func (r *feedbackQueryRepository) FindAll(ctx context.Context, filter *dto.FeedbackFilter) ([]*model.Feedback, error) {
	bsonFilter := bson.M{}
	if filter.Category != "" {
		bsonFilter["category"] = filter.Category
	}
	if filter.Rating > 0 {
		bsonFilter["rating"] = filter.Rating
	}

	cursor, err := r.collection.Find(ctx, bsonFilter)
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