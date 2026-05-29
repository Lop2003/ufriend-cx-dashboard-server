package http

import (
	"github.com/gofiber/fiber/v2"
	commandHandler "ufriend-cx-dashboard-server/internal/follow_up/command/handler"
)

func RegisterFollowUpHTTPRoutes(router fiber.Router, ch *commandHandler.FollowUpCommandHandler) {
	api := router.Group("/api/follow-ups")
	api.Post("/", ch.CreateFollowUp)
	api.Patch("/:id", ch.UpdateStatus)
}
