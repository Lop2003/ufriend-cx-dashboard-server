package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"ufriend-cx-dashboard-server/internal/feedback/command/dto"
	"ufriend-cx-dashboard-server/internal/feedback/command/usecase"
	"ufriend-cx-dashboard-server/pkg/response"
	"ufriend-cx-dashboard-server/pkg/validator"
)

type FeedbackCommandHandler struct {
	usecase usecase.FeedbackCommandUsecase
}

func NewFeedbackCommandHandler(uc usecase.FeedbackCommandUsecase) *FeedbackCommandHandler {
	return &FeedbackCommandHandler{usecase: uc}
}

func (h *FeedbackCommandHandler) CreateFeedback(c *fiber.Ctx) error {
	var req dto.CreateFeedbackRequest
	if err := validator.ValidateBody(c, &req); err != nil {
		return err
	}

	feedback, err := h.usecase.CreateFeedback(c.Context(), &req)
	if err != nil {
		slog.Error("create feedback failed", "error", err)
		return response.ErrorServer(c, "failed to create feedback", err)
	}

	return response.OK(c, "feedback created", feedback)
}