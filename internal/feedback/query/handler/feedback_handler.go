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

func (h *FeedbackQueryHandler) GetFeedbackStats(c *fiber.Ctx) error {
	branch := c.Query("branch")

	stats, err := h.usecase.GetStats(c.Context(), branch)
	if err != nil {
		slog.Error("get feedback stats failed", "error", err)
		return response.ErrorServer(c, "failed to get feedback stats", err)
	}

	return response.OK(c, "feedback stats retrieved", stats)
}