package http

import (
	"github.com/gofiber/fiber/v2"
	queryHandler "ufriend-cx-dashboard-server/internal/customer/query/handler"
)

func RegisterCustomerHTTPRoutes(app *fiber.App, qh *queryHandler.CustomerQueryHandler) {
	api := app.Group("/api/customers")
	api.Get("/", qh.ListCustomers)
	api.Get("/:id", qh.GetCustomer)

	stats := app.Group("/api/stats")
	stats.Get("/summary", qh.GetSummary)
	stats.Get("/by-branch", qh.GetByBranch)
}
