package outbound

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"ufriend-cx-dashboard-server/internal/customer/adapter/inbound"
)

// CustomerAdapter — outbound adapter สำหรับ follow_up domain
// ใช้เรียก customer domain ผ่าน inbound adapter แทนการ import internal ตรง
type CustomerAdapter interface {
	CustomerExists(ctx context.Context, customerID primitive.ObjectID) (bool, error)
}

type customerAdapter struct {
	customer inbound.CustomerAdapter
}

func NewCustomerAdapter(c inbound.CustomerAdapter) CustomerAdapter {
	return &customerAdapter{customer: c}
}

func (a *customerAdapter) CustomerExists(ctx context.Context, customerID primitive.ObjectID) (bool, error) {
	return a.customer.Exists(ctx, customerID)
}
