package usecase

import (
	"context"
	"fmt"
	"time"

	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

func (u *customerQueryUsecase) GetCustomer(ctx context.Context, id string) (*dto.CustomerDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	detail, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get customer: %w", err)
	}

	return detail, nil
}