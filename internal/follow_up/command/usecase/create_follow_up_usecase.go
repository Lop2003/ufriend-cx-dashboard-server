package usecase

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"ufriend-cx-dashboard-server/internal/follow_up/command/dto"
	"ufriend-cx-dashboard-server/internal/follow_up/model"
)

func (u *followUpCommandUsecase) CreateFollowUp(ctx context.Context, req *dto.CreateFollowUpRequest) (*model.FollowUp, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	customerId, err := primitive.ObjectIDFromHex(req.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("invalid customer_id: %w", err)
	}

	followUp := &model.FollowUp{
		Id:         primitive.NewObjectID(),
		CustomerId: customerId,
		Type:       req.Type,
		Note:       req.Note,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	if err := u.repo.Create(ctx, followUp); err != nil {
		return nil, fmt.Errorf("create follow up: %w", err)
	}

	return followUp, nil
}