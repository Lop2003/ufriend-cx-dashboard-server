package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if code >= 500 {
		slog.Error("unhandled error",
			"method", c.Method(),
			"path", c.Path(),
			"status", code,
			"error", err.Error(),
		)
	}

	message := "internal server error"
	if code == fiber.StatusNotFound {
		message = "ไม่พบ API ที่เรียก"
	} else if code == fiber.StatusUnauthorized {
		message = "ไม่ได้รับอนุญาต"
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
		"error":   err.Error(),
	})
}