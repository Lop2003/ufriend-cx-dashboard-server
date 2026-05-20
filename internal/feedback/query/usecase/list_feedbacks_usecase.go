package usecase

import (
	"context"
	"fmt"
	"time"

	"ufriend-cx-dashboard-server/internal/feedback/model"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func (u *feedbackQueryUsecase) ListFeedbacks(ctx context.Context, filter *dto.FeedbackFilter) ([]*model.Feedback, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	feedbacks, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list feedbacks: %w", err)
	}

	return feedbacks, nil
}