package http

import (
	"github.com/gofiber/fiber/v2"
	commandHandler "ufriend-cx-dashboard-server/internal/feedback/command/handler"
	queryHandler "ufriend-cx-dashboard-server/internal/feedback/query/handler"
)

func RegisterFeedbackHTTPRoutes(router fiber.Router, qh *queryHandler.FeedbackQueryHandler, ch *commandHandler.FeedbackCommandHandler) {
	api := router.Group("/api/feedbacks")
	api.Get("/", qh.ListFeedbacks)
	api.Get("/stats", qh.GetFeedbackStats)
	api.Post("/", ch.CreateFeedback)
}
