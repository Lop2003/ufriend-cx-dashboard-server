package usecase

import (
	"context"

	"ufriend-cx-dashboard-server/internal/feedback/model"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
	"ufriend-cx-dashboard-server/internal/feedback/query/repository"
)

type FeedbackQueryUsecase interface {
	ListFeedbacks(ctx context.Context, filter *dto.FeedbackFilter) ([]*model.Feedback, error)
}

type feedbackQueryUsecase struct {
	repo repository.FeedbackQueryRepository
}

func NewFeedbackQueryUsecase(repo repository.FeedbackQueryRepository) FeedbackQueryUsecase {
	return &feedbackQueryUsecase{repo: repo}
}