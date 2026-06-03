package usecase

import (
	"context"

	"ufriend-cx-dashboard-server/internal/feedback/adapter/outbound"
	"ufriend-cx-dashboard-server/internal/feedback/command/dto"
	"ufriend-cx-dashboard-server/internal/feedback/command/repository"
	"ufriend-cx-dashboard-server/internal/feedback/model"
)

type FeedbackCommandUsecase interface {
	CreateFeedback(ctx context.Context, req *dto.CreateFeedbackRequest) (*model.Feedback, error)
}

type feedbackCommandUsecase struct {
	repo            repository.FeedbackCommandRepository
	customerAdapter outbound.CustomerAdapter
}

func NewFeedbackCommandUsecase(repo repository.FeedbackCommandRepository, customerAdapter outbound.CustomerAdapter) FeedbackCommandUsecase {
	return &feedbackCommandUsecase{repo: repo, customerAdapter: customerAdapter}
}