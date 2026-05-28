package usecase

import (
	"context"

	"ufriend-cx-dashboard-server/internal/follow_up/query/dto"
	"ufriend-cx-dashboard-server/internal/follow_up/query/repository"
)

type FollowUpQueryUsecase interface {
	ListFollowUps(ctx context.Context, filter *dto.FollowUpFilter) (*dto.FollowUpListResponse, error)
}

type followUpQueryUsecase struct {
	repo repository.FollowUpQueryRepository
}

func NewFollowUpQueryUsecase(repo repository.FollowUpQueryRepository) FollowUpQueryUsecase {
	return &followUpQueryUsecase{repo: repo}
}
