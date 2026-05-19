package follow_up

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	commandHandler "ufriend-cx-dashboard-server/internal/follow_up/command/handler"
	commandRepository "ufriend-cx-dashboard-server/internal/follow_up/command/repository"
	commandUsecase "ufriend-cx-dashboard-server/internal/follow_up/command/usecase"
)

type FollowUpDomain struct {
	commandHandler *commandHandler.FollowUpCommandHandler
}

func NewFollowUpDomain(db *mongo.Database) *FollowUpDomain {
	commandRepo := commandRepository.NewFollowUpCommandRepository(db)
	commandUC := commandUsecase.NewFollowUpCommandUsecase(commandRepo)
	commandH := commandHandler.NewFollowUpCommandHandler(commandUC)

	return &FollowUpDomain{
		commandHandler: commandH,
	}
}

func (d *FollowUpDomain) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api/follow-ups")
	api.Post("/", d.commandHandler.CreateFollowUp)
	api.Patch("/:id", d.commandHandler.UpdateStatus)
}