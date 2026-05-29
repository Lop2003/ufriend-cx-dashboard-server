package http

import (
	"github.com/gofiber/fiber/v2"
	"ufriend-cx-dashboard-server/internal/auth/handler"
)

func RegisterAuthHTTPRoutes(app *fiber.App, h *handler.AuthHandler) {
	auth := app.Group("/api/auth")
	auth.Get("/lark", h.LarkRedirect)
	auth.Get("/lark/callback", h.LarkCallback)
	auth.Get("/me", h.GetMe)
	auth.Post("/logout", h.Logout)
}
