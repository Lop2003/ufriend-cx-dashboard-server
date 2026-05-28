package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"ufriend-cx-dashboard-server/internal/follow_up/query/dto"
	"ufriend-cx-dashboard-server/internal/follow_up/query/usecase"
	"ufriend-cx-dashboard-server/pkg/response"
)

type FollowUpQueryHandler struct {
	usecase usecase.FollowUpQueryUsecase
}

func NewFollowUpQueryHandler(uc usecase.FollowUpQueryUsecase) *FollowUpQueryHandler {
	return &FollowUpQueryHandler{usecase: uc}
}

func (h *FollowUpQueryHandler) ListFollowUps(c *fiber.Ctx) error {
	var filter dto.FollowUpFilter
	if err := c.QueryParser(&filter); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid filter", "ERR_PARSE")
	}

	res, err := h.usecase.ListFollowUps(c.Context(), &filter)
	if err != nil {
		slog.Error("list follow ups failed", "error", err)
		return response.ErrorServer(c, "failed to list follow ups", err)
	}

	return response.OK(c, "follow ups retrieved", res)
}
