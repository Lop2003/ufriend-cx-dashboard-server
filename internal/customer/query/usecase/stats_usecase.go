package usecase

import (
	"context"
	"fmt"
	"time"

	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

func (u *customerQueryUsecase) GetSummary(ctx context.Context, period string) (*dto.SummaryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	summary, err := u.repo.GetSummary(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}

	return summary, nil
}

func (u *customerQueryUsecase) GetByBranch(ctx context.Context, period string, branch string) ([]*dto.BranchStat, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stats, err := u.repo.GetByBranch(ctx, period, branch)
	if err != nil {
		return nil, fmt.Errorf("get by branch: %w", err)
	}

	return stats, nil
}

func (u *customerQueryUsecase) GetDailyStats(ctx context.Context, period string) ([]dto.DailyStatEntry, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stats, err := u.repo.GetDailyStats(ctx, period)
	if err != nil {
		return nil, fmt.Errorf("get daily stats: %w", err)
	}

	return stats, nil
}