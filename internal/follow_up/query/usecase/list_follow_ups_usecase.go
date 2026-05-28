package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"ufriend-cx-dashboard-server/internal/follow_up/query/dto"
)

func (u *followUpQueryUsecase) ListFollowUps(ctx context.Context, filter *dto.FollowUpFilter) (*dto.FollowUpListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	items, total, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list follow ups: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	return &dto.FollowUpListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
