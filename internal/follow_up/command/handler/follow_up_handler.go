package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"ufriend-cx-dashboard-server/internal/follow_up/command/dto"
	"ufriend-cx-dashboard-server/internal/follow_up/command/usecase"
	"ufriend-cx-dashboard-server/pkg/response"
	"ufriend-cx-dashboard-server/pkg/validator"
)

type FollowUpCommandHandler struct {
	usecase usecase.FollowUpCommandUsecase
}

func NewFollowUpCommandHandler(uc usecase.FollowUpCommandUsecase) *FollowUpCommandHandler {
	return &FollowUpCommandHandler{usecase: uc}
}

func (h *FollowUpCommandHandler) CreateFollowUp(c *fiber.Ctx) error {
	var req dto.CreateFollowUpRequest
	if err := validator.ValidateBody(c, &req); err != nil {
		return err
	}

	followUp, err := h.usecase.CreateFollowUp(c.Context(), &req)
	if err != nil {
		slog.Error("create follow up failed", "error", err)
		return response.ErrorServer(c, "failed to create follow up", err)
	}

	return response.OK(c, "follow up created", followUp)
}

func (h *FollowUpCommandHandler) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.UpdateFollowUpStatusRequest
	if err := validator.ValidateBody(c, &req); err != nil {
		return err
	}

	if err := h.usecase.UpdateStatus(c.Context(), id, &req); err != nil {
		slog.Error("update follow up status failed", "error", err)
		return response.ErrorServer(c, "failed to update status", err)
	}

	return response.OK(c, "follow up status updated", nil)
}