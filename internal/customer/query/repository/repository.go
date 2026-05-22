package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/customer/model"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

type CustomerQueryRepository interface {
	FindAll(ctx context.Context, filter *dto.CustomerFilter) ([]*model.Customer, error)
	FindByID(ctx context.Context, id string) (*dto.CustomerDetailResponse, error)
	GetSummary(ctx context.Context) (*dto.SummaryResponse, error)
	GetByBranch(ctx context.Context) ([]*dto.BranchStat, error)
}

type customerQueryRepository struct {
	db *mongo.Database
}

func NewCustomerQueryRepository(db *mongo.Database) CustomerQueryRepository {
	return &customerQueryRepository{db: db}
}