package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ufriend-cx-dashboard-server/internal/feedback/model"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func (r *feedbackQueryRepository) FindAll(ctx context.Context, filter *dto.FeedbackFilter) ([]*model.Feedback, error) {
	bsonFilter := bson.M{}
	if filter.Branch != "" {
		bsonFilter["branch"] = filter.Branch
	}
	if filter.Category != "" {
		bsonFilter["category"] = filter.Category
	}
	if filter.Rating > 0 {
		bsonFilter["rating"] = filter.Rating
	}

	opts := options.Find()

	// Enforce pagination to prevent out-of-memory crashes on massive data
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Cap maximum limit
	}

	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}}) // Show newest reviews first

	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
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