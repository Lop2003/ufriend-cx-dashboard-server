package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const sessionCookieName = "ufriend_session"

// SessionAuth ตรวจสอบ session cookie ก่อนอนุญาตเข้าถึง protected routes
// ถ้า session ไม่มีหรือหมดอายุ → 401 Unauthorized
func SessionAuth(db *mongo.Database) fiber.Handler {
	col := db.Collection("auth_sessions")

	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies(sessionCookieName)
		if sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "ไม่ได้รับอนุญาต กรุณาเข้าสู่ระบบ",
				"error":   "ERR_UNAUTHORIZED",
			})
		}

		// ตรวจว่ามี session อยู่ใน database จริง
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		var session bson.M
		err := col.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่",
				"error":   "ERR_SESSION_EXPIRED",
			})
		}

		// ตรวจ expires_at ว่าหมดอายุหรือยัง
		if expiresAt, ok := session["expires_at"].(time.Time); ok {
			if time.Now().After(expiresAt) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่",
					"error":   "ERR_SESSION_EXPIRED",
				})
			}
		}

		return c.Next()
	}
}
