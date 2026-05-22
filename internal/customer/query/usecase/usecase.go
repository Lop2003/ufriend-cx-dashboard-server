package usecase

import (
	"context"

	"ufriend-cx-dashboard-server/internal/customer/model"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
	"ufriend-cx-dashboard-server/internal/customer/query/repository"
)

type CustomerQueryUsecase interface {
	ListCustomers(ctx context.Context, filter *dto.CustomerFilter) ([]*model.Customer, error)
	GetCustomer(ctx context.Context, id string) (*dto.CustomerDetailResponse, error)
	GetSummary(ctx context.Context) (*dto.SummaryResponse, error)
	GetByBranch(ctx context.Context) ([]*dto.BranchStat, error)
}

type customerQueryUsecase struct {
	repo repository.CustomerQueryRepository
}

func NewCustomerQueryUsecase(repo repository.CustomerQueryRepository) CustomerQueryUsecase {
	return &customerQueryUsecase{repo: repo}
}