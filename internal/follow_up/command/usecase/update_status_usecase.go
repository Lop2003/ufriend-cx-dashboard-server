package usecase

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"ufriend-cx-dashboard-server/internal/follow_up/command/dto"
)

func (u *followUpCommandUsecase) UpdateStatus(ctx context.Context, id string, req *dto.UpdateFollowUpStatusRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid follow up id: %w", err)
	}

	if err := u.repo.UpdateStatus(ctx, objId, req.Status); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	return nil
}