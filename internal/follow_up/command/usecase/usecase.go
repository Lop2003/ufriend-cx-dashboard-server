package usecase

import (
	"context"

	"ufriend-cx-dashboard-server/internal/follow_up/adapter/outbound"
	"ufriend-cx-dashboard-server/internal/follow_up/command/dto"
	"ufriend-cx-dashboard-server/internal/follow_up/command/repository"
	"ufriend-cx-dashboard-server/internal/follow_up/model"
)

type FollowUpCommandUsecase interface {
	CreateFollowUp(ctx context.Context, req *dto.CreateFollowUpRequest) (*model.FollowUp, error)
	UpdateStatus(ctx context.Context, id string, req *dto.UpdateFollowUpStatusRequest) error
}

type followUpCommandUsecase struct {
	repo            repository.FollowUpCommandRepository
	customerAdapter outbound.CustomerAdapter
}

func NewFollowUpCommandUsecase(repo repository.FollowUpCommandRepository, customerAdapter outbound.CustomerAdapter) FollowUpCommandUsecase {
	return &followUpCommandUsecase{repo: repo, customerAdapter: customerAdapter}
}