package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"ufriend-cx-dashboard-server/internal/customer"
	"ufriend-cx-dashboard-server/internal/feedback"
	"ufriend-cx-dashboard-server/internal/follow_up"
	"ufriend-cx-dashboard-server/pkg/database"
)

func main() {
	// ตั้งค่า MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	db, err := database.NewMongoDB(mongoURI, "ufriend_cx")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.Close()

	// ตั้งค่า Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New())

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// Register domains
	customerDomain := customer.NewCustomerDomain(db.GetDatabase())
	customerDomain.RegisterRoutes(app)

	feedbackDomain := feedback.NewFeedbackDomain(db.GetDatabase())
	feedbackDomain.RegisterRoutes(app)

	followUpDomain := follow_up.NewFollowUpDomain(db.GetDatabase())
	followUpDomain.RegisterRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	log.Printf("Starting server on %s", port)
	log.Fatal(app.Listen(port))
}
