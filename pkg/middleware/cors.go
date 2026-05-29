package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORSConfig สร้าง CORS middleware — รับ allowedOrigins จาก env (comma-separated)
// ถ้าไม่ได้ตั้งค่าจะ fallback เป็น localhost development defaults
func CORSConfig(allowedOrigins string) fiber.Handler {
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:3001"
	}
	return cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,PATCH,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	})
}