package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
	"ufriend-cx-dashboard-server/internal/feedback/query/usecase"
	"ufriend-cx-dashboard-server/pkg/response"
)

type FeedbackQueryHandler struct {
	usecase usecase.FeedbackQueryUsecase
}

func NewFeedbackQueryHandler(uc usecase.FeedbackQueryUsecase) *FeedbackQueryHandler {
	return &FeedbackQueryHandler{usecase: uc}
}

func (h *FeedbackQueryHandler) ListFeedbacks(c *fiber.Ctx) error {
	var filter dto.FeedbackFilter
	if err := c.QueryParser(&filter); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid filter", "ERR_PARSE")
	}

	feedbacks, err := h.usecase.ListFeedbacks(c.Context(), &filter)
	if err != nil {
		slog.Error("list feedbacks failed", "error", err)
		return response.ErrorServer(c, "failed to list feedbacks", err)
	}

	return response.OK(c, "feedbacks retrieved", feedbacks)
}