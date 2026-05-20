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

	slog.Error("unhandled error",
		"method", c.Method(),
		"path", c.Path(),
		"status", code,
		"error", err.Error(),
	)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": "internal server error",
		"error":   err.Error(),
	})
}