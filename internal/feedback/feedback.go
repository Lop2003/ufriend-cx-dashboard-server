package feedback

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	commandHandler "ufriend-cx-dashboard-server/internal/feedback/command/handler"
	commandRepository "ufriend-cx-dashboard-server/internal/feedback/command/repository"
	commandUsecase "ufriend-cx-dashboard-server/internal/feedback/command/usecase"
	queryHandler "ufriend-cx-dashboard-server/internal/feedback/query/handler"
	queryRepository "ufriend-cx-dashboard-server/internal/feedback/query/repository"
	queryUsecase "ufriend-cx-dashboard-server/internal/feedback/query/usecase"
	"ufriend-cx-dashboard-server/internal/feedback/router/http"
)

type FeedbackDomain struct {
	queryHandler   *queryHandler.FeedbackQueryHandler
	commandHandler *commandHandler.FeedbackCommandHandler
}

func NewFeedbackDomain(db *mongo.Database) *FeedbackDomain {
	queryRepo := queryRepository.NewFeedbackQueryRepository(db)
	queryUC := queryUsecase.NewFeedbackQueryUsecase(queryRepo)
	queryH := queryHandler.NewFeedbackQueryHandler(queryUC)

	commandRepo := commandRepository.NewFeedbackCommandRepository(db)
	commandUC := commandUsecase.NewFeedbackCommandUsecase(commandRepo)
	commandH := commandHandler.NewFeedbackCommandHandler(commandUC)

	return &FeedbackDomain{
		queryHandler:   queryH,
		commandHandler: commandH,
	}
}

func (d *FeedbackDomain) RegisterRoutes(app *fiber.App) {
	http.RegisterFeedbackHTTPRoutes(app, d.queryHandler, d.commandHandler)
}