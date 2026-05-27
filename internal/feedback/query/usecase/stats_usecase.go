package usecase

import (
	"context"
	"fmt"
	"time"

	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func (u *feedbackQueryUsecase) GetStats(ctx context.Context, branch string) (*dto.FeedbackStatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stats, err := u.repo.GetStats(ctx, branch)
	if err != nil {
		return nil, fmt.Errorf("get feedback stats usecase: %w", err)
	}

	return stats, nil
}
