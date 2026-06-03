package follow_up

import (
	"github.com/gofiber/fiber/v2"
	customerInbound "ufriend-cx-dashboard-server/internal/customer/adapter/inbound"
	"ufriend-cx-dashboard-server/internal/follow_up/adapter/outbound"
	commandHandler "ufriend-cx-dashboard-server/internal/follow_up/command/handler"
	commandRepository "ufriend-cx-dashboard-server/internal/follow_up/command/repository"
	commandUsecase "ufriend-cx-dashboard-server/internal/follow_up/command/usecase"
	"ufriend-cx-dashboard-server/internal/follow_up/router/http"
	"go.mongodb.org/mongo-driver/mongo"
)

type FollowUpDomain struct {
	commandHandler *commandHandler.FollowUpCommandHandler
}

func NewFollowUpDomain(db *mongo.Database, customerAdapter customerInbound.CustomerAdapter) *FollowUpDomain {
	// 1. Outbound adapters
	custAdapter := outbound.NewCustomerAdapter(customerAdapter)

	// 2. Command side (bottom-up)
	commandRepo := commandRepository.NewFollowUpCommandRepository(db)
	commandUC := commandUsecase.NewFollowUpCommandUsecase(commandRepo, custAdapter)
	commandH := commandHandler.NewFollowUpCommandHandler(commandUC)

	return &FollowUpDomain{
		commandHandler: commandH,
	}
}

func (d *FollowUpDomain) RegisterRoutes(router fiber.Router) {
	http.RegisterFollowUpHTTPRoutes(router, d.commandHandler)
}