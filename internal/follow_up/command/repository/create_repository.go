package repository

import (
	"context"

	"ufriend-cx-dashboard-server/internal/follow_up/model"
)

func (r *followUpCommandRepository) Create(ctx context.Context, followUp *model.FollowUp) error {
	_, err := r.collection.InsertOne(ctx, followUp)
	return err
}