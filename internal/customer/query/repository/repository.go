package repository

import (
	"context"

	"ufriend-cx-dashboard-server/internal/customer/model"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"

	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerQueryRepository interface {
	FindAll(ctx context.Context, filter *dto.CustomerFilter) ([]*model.Customer, int64, error)
	FindByID(ctx context.Context, id string) (*dto.CustomerDetailResponse, error)
	GetSummary(ctx context.Context, period string) (*dto.SummaryResponse, error)
	GetByBranch(ctx context.Context, period string, branch string) ([]*dto.BranchStat, error)
	GetDailyStats(ctx context.Context, period string) ([]dto.DailyStatEntry, error)
}

type customerQueryRepository struct {
	db *mongo.Database
}

func NewCustomerQueryRepository(db *mongo.Database) CustomerQueryRepository {
	return &customerQueryRepository{db: db}
}

