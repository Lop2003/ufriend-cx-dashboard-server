package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/feedback/model"
)

type FeedbackCommandRepository interface {
	Create(ctx context.Context, feedback *model.Feedback) error
}

type feedbackCommandRepository struct {
	collection *mongo.Collection
}

func NewFeedbackCommandRepository(db *mongo.Database) FeedbackCommandRepository {
	return &feedbackCommandRepository{collection: db.Collection("feedbacks")}
}