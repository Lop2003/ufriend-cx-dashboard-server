package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"ufriend-cx-dashboard-server/internal/feedback/command/dto"
	"ufriend-cx-dashboard-server/internal/feedback/command/usecase"
	"ufriend-cx-dashboard-server/pkg/response"
)

type FeedbackCommandHandler struct {
	usecase usecase.FeedbackCommandUsecase
}

func NewFeedbackCommandHandler(uc usecase.FeedbackCommandUsecase) *FeedbackCommandHandler {
	return &FeedbackCommandHandler{usecase: uc}
}

func (h *FeedbackCommandHandler) CreateFeedback(c *fiber.Ctx) error {
	var req dto.CreateFeedbackRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body", "ERR_PARSE")
	}

	feedback, err := h.usecase.CreateFeedback(c.Context(), &req)
	if err != nil {
		slog.Error("create feedback failed", "error", err)
		return response.ErrorServer(c, "failed to create feedback", err)
	}

	return response.OK(c, "feedback created", feedback)
}