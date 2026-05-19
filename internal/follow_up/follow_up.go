package follow_up

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type FollowUpDomain struct {
}

func NewFollowUpDomain(db *mongo.Database) *FollowUpDomain {
	return &FollowUpDomain{}
}

func (d *FollowUpDomain) RegisterRoutes(app *fiber.App) {
}
