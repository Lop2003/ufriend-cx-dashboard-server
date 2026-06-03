package repository

import (
	"context"
	"fmt"

	"ufriend-cx-dashboard-server/internal/auth/model"
	"ufriend-cx-dashboard-server/pkg/crypto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateTokens atomically updates tokens only if tokenVersion matches the current value in DB.
// Uses token_version (integer) for optimistic locking instead of matching on encrypted tokens,
// because AES-GCM uses random nonce → same plaintext produces different ciphertext each time.
// Returns (true, nil) if update succeeded, (false, nil) if another request already refreshed.
func (r *authRepository) UpdateTokens(ctx context.Context, sessionID string, tokenVersion int, session *model.AuthSession) (bool, error) {
	// Encrypt new tokens ก่อนบันทึก
	encAccessToken, err := crypto.Encrypt(session.AccessToken, r.encryptionKey)
	if err != nil {
		return false, fmt.Errorf("encrypt access_token: %w", err)
	}
	encRefreshToken, err := crypto.Encrypt(session.RefreshToken, r.encryptionKey)
	if err != nil {
		return false, fmt.Errorf("encrypt refresh_token: %w", err)
	}

	// Atomic update: match session_id + token_version
	// ถ้าคนอื่น refresh ไปแล้ว → token_version จะเพิ่มขึ้น → match ไม่ได้ → MatchedCount = 0
	filter := bson.M{
		"session_id":    sessionID,
		"token_version": tokenVersion,
	}
	update := bson.M{
		"$set": bson.M{
			"access_token":  encAccessToken,
			"refresh_token": encRefreshToken,
			"expires_at":    session.ExpiresAt,
		},
		"$inc": bson.M{
			"token_version": 1,
		},
	}

	opts := options.Update()
	result, err := r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return false, fmt.Errorf("update tokens: %w", err)
	}

	return result.MatchedCount > 0, nil
}
