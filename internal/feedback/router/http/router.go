package http

import (
	"github.com/gofiber/fiber/v2"
	"ufriend-cx-dashboard-server/internal/feedback/command/handler"
	queryHandler "ufriend-cx-dashboard-server/internal/feedback/query/handler"
)

func RegisterFeedbackHTTPRoutes(app *fiber.App, qh *queryHandler.FeedbackQueryHandler, ch *handler.FeedbackCommandHandler) {
	api := app.Group("/api/feedbacks")
	api.Get("/", qh.List)
	api.Post("/", ch.Create)
}
