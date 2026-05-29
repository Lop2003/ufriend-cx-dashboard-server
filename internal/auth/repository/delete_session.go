package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (r *authRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"session_id": sessionID})
	return err
}
