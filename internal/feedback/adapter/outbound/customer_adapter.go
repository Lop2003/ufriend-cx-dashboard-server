package outbound

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"ufriend-cx-dashboard-server/internal/customer/adapter/inbound"
)

// CustomerAdapter — outbound adapter สำหรับ feedback domain
// ใช้เรียก customer domain ผ่าน inbound adapter แทนการ import internal ตรง
type CustomerAdapter interface {
	GetBranchByCustomerId(ctx context.Context, customerID primitive.ObjectID) (string, error)
}

type customerAdapter struct {
	customer inbound.CustomerAdapter
}

func NewCustomerAdapter(c inbound.CustomerAdapter) CustomerAdapter {
	return &customerAdapter{customer: c}
}

func (a *customerAdapter) GetBranchByCustomerId(ctx context.Context, customerID primitive.ObjectID) (string, error) {
	return a.customer.GetBranch(ctx, customerID)
}
