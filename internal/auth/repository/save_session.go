package repository

import (
	"context"
	"fmt"

	"ufriend-cx-dashboard-server/internal/auth/model"
	"ufriend-cx-dashboard-server/pkg/crypto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *authRepository) Save(ctx context.Context, session *model.AuthSession) error {
	// Encrypt tokens ก่อนบันทึกลง database เพื่อป้องกัน token leak
	encAccessToken, err := crypto.Encrypt(session.AccessToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("encrypt access_token: %w", err)
	}
	encRefreshToken, err := crypto.Encrypt(session.RefreshToken, r.encryptionKey)
	if err != nil {
		return fmt.Errorf("encrypt refresh_token: %w", err)
	}

	filter := bson.M{"session_id": session.SessionID}
	update := bson.M{
		"$set": bson.M{
			"user_id":       session.UserID,
			"access_token":  encAccessToken,
			"refresh_token": encRefreshToken,
			"expires_at":    session.ExpiresAt,
			"token_version": session.TokenVersion,
		},
		"$setOnInsert": bson.M{
			"session_id": session.SessionID,
			"created_at": session.CreatedAt,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

