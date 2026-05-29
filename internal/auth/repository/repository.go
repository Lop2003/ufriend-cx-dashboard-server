package repository

import (
	"context"

	"ufriend-cx-dashboard-server/internal/auth/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	FindBySessionID(ctx context.Context, sessionID string) (*model.AuthSession, error)
	Save(ctx context.Context, session *model.AuthSession) error
	DeleteBySessionID(ctx context.Context, sessionID string) error
}

type authRepository struct {
	collection *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return &authRepository{
		collection: db.Collection("auth_sessions"),
	}
}
