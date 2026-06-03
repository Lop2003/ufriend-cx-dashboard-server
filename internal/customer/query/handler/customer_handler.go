package handler

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
	"ufriend-cx-dashboard-server/internal/customer/query/usecase"
	"ufriend-cx-dashboard-server/pkg/response"
)

type CustomerQueryHandler struct {
	usecase usecase.CustomerQueryUsecase
}

func NewCustomerQueryHandler(uc usecase.CustomerQueryUsecase) *CustomerQueryHandler {
	return &CustomerQueryHandler{usecase: uc}
}

func (h *CustomerQueryHandler) ListCustomers(c *fiber.Ctx) error {
	var filter dto.CustomerFilter
	if err := c.QueryParser(&filter); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid filter", "ERR_PARSE")
	}

	customers, err := h.usecase.ListCustomers(c.Context(), &filter)
	if err != nil {
		slog.Error("list customers failed", "error", err)
		return response.ErrorServer(c, "failed to list customers", err)
	}

	return response.OK(c, "customers retrieved", customers)
}

func (h *CustomerQueryHandler) GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")

	detail, err := h.usecase.GetCustomer(c.Context(), id)
	if err != nil {
		slog.Error("get customer failed", "error", err)
		// แยก error: not found → 404, อื่นๆ (DB down, invalid ID) → 500
		if errors.Is(err, mongo.ErrNoDocuments) {
			return response.Error(c, fiber.StatusNotFound, "customer not found", "ERR_NOT_FOUND")
		}
		return response.ErrorServer(c, "failed to get customer", err)
	}

	return response.OK(c, "customer retrieved", detail)
}

func (h *CustomerQueryHandler) GetSummary(c *fiber.Ctx) error {
	period := c.Query("period")

	summary, err := h.usecase.GetSummary(c.Context(), period)
	if err != nil {
		slog.Error("get summary failed", "error", err)
		return response.ErrorServer(c, "failed to get summary", err)
	}

	return response.OK(c, "summary retrieved", summary)
}

func (h *CustomerQueryHandler) GetByBranch(c *fiber.Ctx) error {
	branch := c.Query("branch")
	period := c.Query("period")

	// ส่ง branch ไป usecase/repo ให้ filter ที่ DB level แทน in-memory
	stats, err := h.usecase.GetByBranch(c.Context(), period, branch)
	if err != nil {
		slog.Error("get by branch failed", "error", err)
		return response.ErrorServer(c, "failed to get branch stats", err)
	}

	return response.OK(c, "branch stats retrieved", stats)
}

func (h *CustomerQueryHandler) GetDailyStats(c *fiber.Ctx) error {
	period := c.Query("period")

	stats, err := h.usecase.GetDailyStats(c.Context(), period)
	if err != nil {
		slog.Error("get daily stats failed", "error", err)
		return response.ErrorServer(c, "failed to get daily stats", err)
	}

	return response.OK(c, "daily stats retrieved", stats)
}