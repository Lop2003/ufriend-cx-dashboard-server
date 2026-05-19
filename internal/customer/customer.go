package customer

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerDomain struct {
}

func NewCustomerDomain(db *mongo.Database) *CustomerDomain {
	return &CustomerDomain{}
}

func (d *CustomerDomain) RegisterRoutes(app *fiber.App) {
}
