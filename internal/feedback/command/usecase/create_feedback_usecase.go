package usecase

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"ufriend-cx-dashboard-server/internal/feedback/command/dto"
	"ufriend-cx-dashboard-server/internal/feedback/logic"
	"ufriend-cx-dashboard-server/internal/feedback/model"
)

func (u *feedbackCommandUsecase) CreateFeedback(ctx context.Context, req *dto.CreateFeedbackRequest) (*model.Feedback, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	customerId, err := primitive.ObjectIDFromHex(req.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("invalid customer_id: %w", err)
	}

	// ดึง branch จาก customer domain ผ่าน outbound adapter (ไม่เรียก customers collection ตรง)
	branch, err := u.customerAdapter.GetBranchByCustomerId(ctx, customerId)
	if err != nil {
		return nil, fmt.Errorf("get customer branch: %w", err)
	}

	feedback := &model.Feedback{
		Id:         primitive.NewObjectID(),
		CustomerId: customerId,
		Rating:     req.Rating,
		Comment:    req.Comment,
		Category:   req.Category,
		Sentiment:  logic.SentimentFromRating(req.Rating),
		Branch:     branch,
		CreatedAt:  time.Now(),
	}

	if err := u.repo.Create(ctx, feedback); err != nil {
		return nil, fmt.Errorf("create feedback: %w", err)
	}

	return feedback, nil
}