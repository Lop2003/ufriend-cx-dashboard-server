package repository

import (
	"context"

	"ufriend-cx-dashboard-server/internal/feedback/model"
)

func (r *feedbackCommandRepository) Create(ctx context.Context, feedback *model.Feedback) error {
	_, err := r.collection.InsertOne(ctx, feedback)
	return err
}