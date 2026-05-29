package model

import "time"

// AuthSession เก็บ token ของผู้ใช้ที่ login ผ่าน Lark OAuth
type AuthSession struct {
	SessionID    string    `bson:"session_id"`
	UserID       string    `bson:"user_id"`
	AccessToken  string    `bson:"access_token"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiresAt    time.Time `bson:"expires_at"`
	CreatedAt    time.Time `bson:"created_at"`
}
