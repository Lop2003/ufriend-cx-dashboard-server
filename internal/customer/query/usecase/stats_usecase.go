package usecase

import (
	"context"
	"fmt"
	"time"

	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

func (u *customerQueryUsecase) GetSummary(ctx context.Context) (*dto.SummaryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	summary, err := u.repo.GetSummary(ctx)
	if err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}

	return summary, nil
}

func (u *customerQueryUsecase) GetByBranch(ctx context.Context) ([]*dto.BranchStat, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	stats, err := u.repo.GetByBranch(ctx)
	if err != nil {
		return nil, fmt.Errorf("get by branch: %w", err)
	}

	return stats, nil
}