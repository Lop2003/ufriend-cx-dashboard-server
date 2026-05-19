package usecase

import (
	"context"
	"fmt"
	"time"

	"ufriend-cx-dashboard-server/internal/customer/model"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

func (u *customerQueryUsecase) ListCustomers(ctx context.Context, filter *dto.CustomerFilter) ([]*model.Customer, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	customers, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}

	return customers, nil
}