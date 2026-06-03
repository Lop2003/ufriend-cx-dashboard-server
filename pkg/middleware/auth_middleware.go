package middleware

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/auth/model"
)

const sessionCookieName = "ufriend_session"

// SessionAuth ตรวจสอบ session cookie ก่อนอนุญาตเข้าถึง protected routes
// ถ้า session ไม่มีหรือหมดอายุ → 401 Unauthorized
// Inject user_id และ session_id ลง c.Locals() เพื่อให้ handler ใช้ audit trail ได้
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

		// Decode ลง struct เพื่อให้ expires_at เป็น time.Time ที่ถูกต้อง
		// (bson.M จะ decode เป็น primitive.DateTime ซึ่ง type assert เป็น time.Time ไม่ได้)
		var session model.AuthSession
		err := col.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่",
				"error":   "ERR_SESSION_EXPIRED",
			})
		}

		// ตรวจ expires_at ว่าหมดอายุหรือยัง
		if time.Now().After(session.ExpiresAt) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่",
				"error":   "ERR_SESSION_EXPIRED",
			})
		}

		// Inject user identity ลง context เพื่อให้ handler ใช้ได้ (audit trail)
		c.Locals("user_id", session.UserID)
		c.Locals("session_id", session.SessionID)

		return c.Next()
	}
}
