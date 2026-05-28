package http

import (
	"github.com/gofiber/fiber/v2"
	commandHandler "ufriend-cx-dashboard-server/internal/follow_up/command/handler"
	queryHandler "ufriend-cx-dashboard-server/internal/follow_up/query/handler"
)

func RegisterFollowUpHTTPRoutes(app *fiber.App, ch *commandHandler.FollowUpCommandHandler, qh *queryHandler.FollowUpQueryHandler) {
	api := app.Group("/api/follow-ups")
	api.Get("/", qh.ListFollowUps)
	api.Post("/", ch.CreateFollowUp)
	api.Patch("/:id", ch.UpdateStatus)
}
