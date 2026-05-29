package repository

import (
	"context"

	"ufriend-cx-dashboard-server/internal/auth/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *authRepository) Save(ctx context.Context, session *model.AuthSession) error {
	filter := bson.M{"session_id": session.SessionID}
	update := bson.M{
		"$set": bson.M{
			"user_id":       session.UserID,
			"access_token":  session.AccessToken,
			"refresh_token": session.RefreshToken,
			"expires_at":    session.ExpiresAt,
		},
		"$setOnInsert": bson.M{
			"session_id": session.SessionID,
			"created_at": session.CreatedAt,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
