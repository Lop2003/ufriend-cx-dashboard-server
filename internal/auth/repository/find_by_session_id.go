package repository

import (
	"context"
	"errors"
	"fmt"

	"ufriend-cx-dashboard-server/internal/auth/model"
	"ufriend-cx-dashboard-server/pkg/crypto"

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

	// Decrypt tokens หลังอ่านจาก database
	session.AccessToken, err = crypto.Decrypt(session.AccessToken, r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt access_token: %w", err)
	}
	session.RefreshToken, err = crypto.Decrypt(session.RefreshToken, r.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("decrypt refresh_token: %w", err)
	}

	return &session, nil
}

