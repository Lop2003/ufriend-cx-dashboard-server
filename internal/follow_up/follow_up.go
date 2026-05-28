package follow_up

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	commandHandler "ufriend-cx-dashboard-server/internal/follow_up/command/handler"
	commandRepository "ufriend-cx-dashboard-server/internal/follow_up/command/repository"
	commandUsecase "ufriend-cx-dashboard-server/internal/follow_up/command/usecase"
	queryHandler "ufriend-cx-dashboard-server/internal/follow_up/query/handler"
	queryRepository "ufriend-cx-dashboard-server/internal/follow_up/query/repository"
	queryUsecase "ufriend-cx-dashboard-server/internal/follow_up/query/usecase"
	"ufriend-cx-dashboard-server/internal/follow_up/router/http"
)

type FollowUpDomain struct {
	commandHandler *commandHandler.FollowUpCommandHandler
	queryHandler   *queryHandler.FollowUpQueryHandler
}

func NewFollowUpDomain(db *mongo.Database) *FollowUpDomain {
	commandRepo := commandRepository.NewFollowUpCommandRepository(db)
	commandUC := commandUsecase.NewFollowUpCommandUsecase(commandRepo)
	commandH := commandHandler.NewFollowUpCommandHandler(commandUC)

	queryRepo := queryRepository.NewFollowUpQueryRepository(db)
	queryUC := queryUsecase.NewFollowUpQueryUsecase(queryRepo)
	queryH := queryHandler.NewFollowUpQueryHandler(queryUC)

	return &FollowUpDomain{
		commandHandler: commandH,
		queryHandler:   queryH,
	}
}

func (d *FollowUpDomain) RegisterRoutes(app *fiber.App) {
	http.RegisterFollowUpHTTPRoutes(app, d.commandHandler, d.queryHandler)
}