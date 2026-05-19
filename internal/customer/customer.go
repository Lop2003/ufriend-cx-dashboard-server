package customer

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	queryHandler "ufriend-cx-dashboard-server/internal/customer/query/handler"
	queryRepository "ufriend-cx-dashboard-server/internal/customer/query/repository"
	queryUsecase "ufriend-cx-dashboard-server/internal/customer/query/usecase"
)

type CustomerDomain struct {
	queryHandler *queryHandler.CustomerQueryHandler
}

func NewCustomerDomain(db *mongo.Database) *CustomerDomain {
	queryRepo := queryRepository.NewCustomerQueryRepository(db)
	queryUC := queryUsecase.NewCustomerQueryUsecase(queryRepo)
	queryH := queryHandler.NewCustomerQueryHandler(queryUC)

	return &CustomerDomain{
		queryHandler: queryH,
	}
}

func (d *CustomerDomain) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/customers")
	api.Get("/", d.queryHandler.ListCustomers)
	api.Get("/:id", d.queryHandler.GetCustomer)

	stats := app.Group("/api/stats")
	stats.Get("/summary", d.queryHandler.GetSummary)
	stats.Get("/by-branch", d.queryHandler.GetByBranch)
}