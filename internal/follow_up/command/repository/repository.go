package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/follow_up/model"
)

type FollowUpCommandRepository interface {
	Create(ctx context.Context, followUp *model.FollowUp) error
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error
}

type followUpCommandRepository struct {
	collection *mongo.Collection
}

func NewFollowUpCommandRepository(db *mongo.Database) FollowUpCommandRepository {
	return &followUpCommandRepository{collection: db.Collection("follow_ups")}
}