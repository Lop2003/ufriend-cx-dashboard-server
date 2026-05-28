package repository

import (
	"context"

	"ufriend-cx-dashboard-server/internal/feedback/model"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"

	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackQueryRepository interface {
	FindAll(ctx context.Context, filter *dto.FeedbackFilter) ([]*model.Feedback, error)
	GetStats(ctx context.Context, branch string) (*dto.FeedbackStatsResponse, error)
}

type feedbackQueryRepository struct {
	collection *mongo.Collection
}

func NewFeedbackQueryRepository(db *mongo.Database) FeedbackQueryRepository {
	return &feedbackQueryRepository{collection: db.Collection("feedbacks")}
}
