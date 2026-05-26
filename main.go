package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"ufriend-cx-dashboard-server/internal/customer"
	"ufriend-cx-dashboard-server/internal/feedback"
	"ufriend-cx-dashboard-server/internal/follow_up"
	"ufriend-cx-dashboard-server/pkg/database"
	"ufriend-cx-dashboard-server/pkg/middleware"
)

func main() {
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	db, err := database.NewMongoDB(mongoURI, "ufriend_cx")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Middleware
	app.Use(middleware.RecoverConfig())
	app.Use(middleware.CORSConfig())

	// Health check
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