package feedback

import (
	"github.com/gofiber/fiber/v2"
	customerInbound "ufriend-cx-dashboard-server/internal/customer/adapter/inbound"
	"ufriend-cx-dashboard-server/internal/feedback/adapter/outbound"
	commandHandler "ufriend-cx-dashboard-server/internal/feedback/command/handler"
	commandRepository "ufriend-cx-dashboard-server/internal/feedback/command/repository"
	commandUsecase "ufriend-cx-dashboard-server/internal/feedback/command/usecase"
	queryHandler "ufriend-cx-dashboard-server/internal/feedback/query/handler"
	queryRepository "ufriend-cx-dashboard-server/internal/feedback/query/repository"
	queryUsecase "ufriend-cx-dashboard-server/internal/feedback/query/usecase"
	"ufriend-cx-dashboard-server/internal/feedback/router/http"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackDomain struct {
	queryHandler   *queryHandler.FeedbackQueryHandler
	commandHandler *commandHandler.FeedbackCommandHandler
}

func NewFeedbackDomain(db *mongo.Database, customerAdapter customerInbound.CustomerAdapter) *FeedbackDomain {
	// 1. Query side (bottom-up)
	queryRepo := queryRepository.NewFeedbackQueryRepository(db)
	queryUC := queryUsecase.NewFeedbackQueryUsecase(queryRepo)
	queryH := queryHandler.NewFeedbackQueryHandler(queryUC)

	// 2. Outbound adapters
	custAdapter := outbound.NewCustomerAdapter(customerAdapter)

	// 3. Command side (bottom-up)
	commandRepo := commandRepository.NewFeedbackCommandRepository(db)
	commandUC := commandUsecase.NewFeedbackCommandUsecase(commandRepo, custAdapter)
	commandH := commandHandler.NewFeedbackCommandHandler(commandUC)

	return &FeedbackDomain{
		queryHandler:   queryH,
		commandHandler: commandH,
	}
}

func (d *FeedbackDomain) RegisterRoutes(router fiber.Router) {
	http.RegisterFeedbackHTTPRoutes(router, d.queryHandler, d.commandHandler)
}