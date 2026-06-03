package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"ufriend-cx-dashboard-server/internal/auth"
	"ufriend-cx-dashboard-server/internal/customer"
	customerInbound "ufriend-cx-dashboard-server/internal/customer/adapter/inbound"
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

	// Request logging — access log สำหรับ debug production
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// CORS — warn ถ้า production ไม่ได้ตั้ง CORS_ORIGINS
	corsOrigins := os.Getenv("CORS_ORIGINS")
	if corsOrigins == "" && os.Getenv("APP_ENV") == "production" {
		slog.Warn("CORS_ORIGINS not set in production — defaulting to localhost origins. Browser requests from production domain will be blocked!")
	}
	app.Use(middleware.CORSConfig(corsOrigins))

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

	authEnabled := larkAppID != "" && larkAppSecret != "" && larkRedirectURI != ""

	if authEnabled {
		authDomain := auth.NewAuthDomain(db.GetDatabase(), auth.Config{
			AppID:           larkAppID,
			AppSecret:       larkAppSecret,
			RedirectURI:     larkRedirectURI,
			ClientURL:       clientURL,
			LarkBaseURL:     larkBaseURL,
			LarkAccountsURL: larkAccountsURL,
			LarkOAuthScope:  larkOAuthScope,
			EncryptionKey:   os.Getenv("SESSION_ENCRYPTION_KEY"),
		})
		authDomain.RegisterRoutes(app)
	} else {
		log.Println("⚠️ Lark OAuth disabled: set LARK_APP_ID, LARK_APP_SECRET, LARK_REDIRECT_URI")
	}

	// Protected API group — ใช้ auth middleware เฉพาะเมื่อ Lark OAuth enabled
	// เพื่อไม่ block dev workflow ที่ไม่ได้ setup Lark
	protected := app.Group("")
	if authEnabled {
		protected.Use(middleware.SessionAuth(db.GetDatabase()))
	}

	// Register domains ภายใต้ protected group
	customerDomain := customer.NewCustomerDomain(db.GetDatabase())
	customerDomain.RegisterRoutes(protected)

	customerAdapter := customerInbound.NewCustomerAdapter(db.GetDatabase())

	feedbackDomain := feedback.NewFeedbackDomain(db.GetDatabase(), customerAdapter)
	feedbackDomain.RegisterRoutes(protected)

	followUpDomain := follow_up.NewFollowUpDomain(db.GetDatabase(), customerAdapter)
	followUpDomain.RegisterRoutes(protected)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	} else {
		port = ":" + port
	}

	// Graceful shutdown — รอ in-flight requests จบก่อนปิด server
	// ป้องกัน request ถูกตัดกลางทางเมื่อ deploy ใหม่
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", port)
		if err := app.Listen(port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	slog.Info("shutting down server gracefully...")

	if err := app.Shutdown(); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("server stopped")
}