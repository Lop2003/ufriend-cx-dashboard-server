package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"ufriend-cx-dashboard-server/internal/auth"
	"ufriend-cx-dashboard-server/internal/customer"
	"ufriend-cx-dashboard-server/internal/feedback"
	"ufriend-cx-dashboard-server/internal/follow_up"
	"ufriend-cx-dashboard-server/pkg/database"
	"ufriend-cx-dashboard-server/pkg/middleware"
)

func main() {
	// โหลด .env จากโฟลเดอร์ server และ root โปรเจกต์ (รองรับการรันจาก directory ใดก็ได้)
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "ufriend_cx"
	}

	db, err := database.NewMongoDB(mongoURI, dbName)
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

	// Lark OAuth env
	larkAppID := os.Getenv("LARK_APP_ID")
	larkAppSecret := os.Getenv("LARK_APP_SECRET")
	larkRedirectURI := os.Getenv("LARK_REDIRECT_URI")
	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:5173"
	}
	larkBaseURL := os.Getenv("LARK_BASE_URL")
	larkAccountsURL := os.Getenv("LARK_ACCOUNTS_URL")
	larkOAuthScope := os.Getenv("LARK_OAUTH_SCOPE")

	if larkAppID != "" && larkAppSecret != "" && larkRedirectURI != "" {
		authDomain := auth.NewAuthDomain(db.GetDatabase(), auth.Config{
			AppID:           larkAppID,
			AppSecret:       larkAppSecret,
			RedirectURI:     larkRedirectURI,
			ClientURL:       clientURL,
			LarkBaseURL:     larkBaseURL,
			LarkAccountsURL: larkAccountsURL,
			LarkOAuthScope:  larkOAuthScope,
		})
		authDomain.RegisterRoutes(app)
	} else {
		log.Println("⚠️ Lark OAuth disabled: set LARK_APP_ID, LARK_APP_SECRET, LARK_REDIRECT_URI")
	}

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