package usecase

import (
	"context"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func (u *feedbackQueryUsecase) GetFeedbacksByCustomerID(ctx context.Context, customerId string) ([]*dto.FeedbackResponse, error) {
	feedbacks, err := u.repo.FindByCustomerID(ctx, customerId)
	if err != nil {
		return nil, err
	}
	return mapToResponseList(feedbacks), nil
}
