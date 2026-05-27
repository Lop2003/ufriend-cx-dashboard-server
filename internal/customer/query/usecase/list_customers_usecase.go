package usecase

import (
	"context"
	"fmt"
	"math"
	"time"

	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

func (u *customerQueryUsecase) ListCustomers(ctx context.Context, filter *dto.CustomerFilter) (*dto.PaginatedCustomerResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	customers, total, err := u.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list customers: %w", err)
	}

	// Normalise pagination values (ตรงกับที่ repository ใช้จริง)
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		// ไม่ได้ระบุ limit → return ทั้งหมด
		limit = int(total)
		if limit < 1 {
			limit = 1
		}
	} else if limit > 100 {
		limit = 100
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	// Map model → DTO
	items := make([]*dto.CustomerResponse, 0, len(customers))
	for _, c := range customers {
		items = append(items, &dto.CustomerResponse{
			Id:         c.Id.Hex(),
			Name:       c.Name,
			Phone:      c.Phone,
			Product:    c.Product,
			Branch:     c.Branch,
			PlanMonths: c.PlanMonths,
			Status:     c.Status,
			CreatedAt:  c.CreatedAt,
		})
	}

	return &dto.PaginatedCustomerResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
