package usecase

import (
	"context"
	"ufriend-cx-dashboard-server/internal/feedback/model"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func mapToResponseList(feedbacks []*model.Feedback) []*dto.FeedbackResponse {
	var list []*dto.FeedbackResponse
	for _, f := range feedbacks {
		list = append(list, &dto.FeedbackResponse{
			Id:         f.Id.Hex(),
			CustomerId: f.CustomerId.Hex(),
			Rating:     f.Rating,
			Comment:    f.Comment,
			Category:   f.Category,
			Sentiment:  f.Sentiment,
			CreatedAt:  f.CreatedAt,
		})
	}
	if list == nil {
		list = make([]*dto.FeedbackResponse, 0)
	}
	return list
}

func (u *feedbackQueryUsecase) GetFeedbacks(ctx context.Context, filter *dto.FeedbackFilter) ([]*dto.FeedbackResponse, error) {
	feedbacks, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}
	return mapToResponseList(feedbacks), nil
}
