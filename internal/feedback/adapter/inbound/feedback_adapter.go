package inbound

import (
	"context"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
	"ufriend-cx-dashboard-server/internal/feedback/query/usecase"
)

type FeedbackAdapter interface {
	GetFeedbacksByCustomerID(ctx context.Context, customerId string) ([]*dto.FeedbackResponse, error)
}

type feedbackAdapter struct {
	queryUsecase usecase.FeedbackQueryUsecase
}

func NewFeedbackAdapter(qu usecase.FeedbackQueryUsecase) FeedbackAdapter {
	return &feedbackAdapter{queryUsecase: qu}
}

func (a *feedbackAdapter) GetFeedbacksByCustomerID(ctx context.Context, customerId string) ([]*dto.FeedbackResponse, error) {
	return a.queryUsecase.GetFeedbacksByCustomerID(ctx, customerId)
}
