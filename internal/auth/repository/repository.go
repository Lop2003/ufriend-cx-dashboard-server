package repository

import (
	"context"

	"ufriend-cx-dashboard-server/internal/auth/model"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	FindBySessionID(ctx context.Context, sessionID string) (*model.AuthSession, error)
	Save(ctx context.Context, session *model.AuthSession) error
	UpdateTokens(ctx context.Context, sessionID string, tokenVersion int, session *model.AuthSession) (bool, error)
	DeleteBySessionID(ctx context.Context, sessionID string) error
}

type authRepository struct {
	collection    *mongo.Collection
	encryptionKey string
}

func NewAuthRepository(db *mongo.Database, encryptionKey string) AuthRepository {
	return &authRepository{
		collection:    db.Collection("auth_sessions"),
		encryptionKey: encryptionKey,
	}
}
