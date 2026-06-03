package handler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
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
		return err // ValidateBody ส่ง response กลับเองแล้ว
	}

	feedback, err := h.usecase.CreateFeedback(c.Context(), &req)
	if err != nil {
		slog.Error("create feedback failed", "error", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return response.Error(c, fiber.StatusNotFound, "customer not found", "ERR_NOT_FOUND")
		}
		return response.ErrorServer(c, "failed to create feedback", err)
	}

	return response.OK(c, "feedback created", feedback)
}