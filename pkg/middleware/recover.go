package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func RecoverConfig() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			slog.Error("panic recovered",
				"method", c.Method(),
				"path", c.Path(),
				"error", e,
			)
		},
	})
}