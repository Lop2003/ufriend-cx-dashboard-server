package feedback

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackDomain struct {
}

func NewFeedbackDomain(db *mongo.Database) *FeedbackDomain {
	return &FeedbackDomain{}
}

func (d *FeedbackDomain) RegisterRoutes(app *fiber.App) {
}
