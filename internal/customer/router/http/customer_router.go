package http

import (
	"github.com/gofiber/fiber/v2"
	queryHandler "ufriend-cx-dashboard-server/internal/customer/query/handler"
)

func RegisterCustomerHTTPRoutes(router fiber.Router, qh *queryHandler.CustomerQueryHandler) {
	api := router.Group("/api/customers")
	api.Get("/", qh.ListCustomers)
	api.Get("/:id", qh.GetCustomer)

	stats := router.Group("/api/stats")
	stats.Get("/summary", qh.GetSummary)
	stats.Get("/branches", qh.GetByBranch)
}

