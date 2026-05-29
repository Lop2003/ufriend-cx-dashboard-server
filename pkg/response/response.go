package response

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(APIResponse{Success: true, Message: message, Data: data})
}

func Error(c *fiber.Ctx, statusCode int, message string, errCode string) error {
	return c.Status(statusCode).JSON(APIResponse{Success: false, Message: message, Error: errCode})
}

// ErrorServer ส่ง 500 response — แสดง error detail เฉพาะ development เท่านั้น
// ใช้ APP_ENV env var แทนการเช็ค IP เพื่อให้ปลอดภัยเมื่ออยู่หลัง reverse proxy
func ErrorServer(c *fiber.Ctx, message string, err error) error {
	errStr := "internal server error"
	if os.Getenv("APP_ENV") != "production" {
		errStr = err.Error()
	}
	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{Success: false, Message: message, Error: errStr})
}