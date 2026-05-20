package usecase

import (
	"context"

	"ufriend-cx-dashboard-server/internal/follow_up/command/dto"
	"ufriend-cx-dashboard-server/internal/follow_up/command/repository"
	"ufriend-cx-dashboard-server/internal/follow_up/model"
)

type FollowUpCommandUsecase interface {
	CreateFollowUp(ctx context.Context, req *dto.CreateFollowUpRequest) (*model.FollowUp, error)
	UpdateStatus(ctx context.Context, id string, req *dto.UpdateFollowUpStatusRequest) error
}

type followUpCommandUsecase struct {
	repo repository.FollowUpCommandRepository
}

func NewFollowUpCommandUsecase(repo repository.FollowUpCommandRepository) FollowUpCommandUsecase {
	return &followUpCommandUsecase{repo: repo}
}