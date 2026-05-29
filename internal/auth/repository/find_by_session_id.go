package repository

import (
	"context"
	"errors"

	"ufriend-cx-dashboard-server/internal/auth/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *authRepository) FindBySessionID(ctx context.Context, sessionID string) (*model.AuthSession, error) {
	var session model.AuthSession
	err := r.collection.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}
	return &session, nil
}
