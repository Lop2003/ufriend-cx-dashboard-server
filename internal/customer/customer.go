package customer

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	queryHandler "ufriend-cx-dashboard-server/internal/customer/query/handler"
	queryRepository "ufriend-cx-dashboard-server/internal/customer/query/repository"
	queryUsecase "ufriend-cx-dashboard-server/internal/customer/query/usecase"
	"ufriend-cx-dashboard-server/internal/customer/router/http"
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
	http.RegisterCustomerHTTPRoutes(app, d.queryHandler)
}