package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ufriend-cx-dashboard-server/internal/auth/handler"
	"ufriend-cx-dashboard-server/internal/auth/lark"
	authHTTP "ufriend-cx-dashboard-server/internal/auth/router/http"
	"ufriend-cx-dashboard-server/internal/auth/repository"
	"ufriend-cx-dashboard-server/internal/auth/usecase"
)

type AuthDomain struct {
	handler *handler.AuthHandler
}

type Config struct {
	AppID          string
	AppSecret      string
	RedirectURI    string
	ClientURL      string
	LarkBaseURL    string
	LarkAccountsURL string
	LarkOAuthScope string
}

func NewAuthDomain(db *mongo.Database, cfg Config) *AuthDomain {
	// สร้าง index สำหรับค้นหา session ด้วย session_id
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		col := db.Collection("auth_sessions")
		_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    bson.D{{Key: "session_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		})
		if err != nil {
			slog.Error("failed to create auth_sessions index",
				"collection", "auth_sessions",
				"error", err,
			)
		}
	}()

	larkClient := lark.NewClient(cfg.AppID, cfg.AppSecret, cfg.LarkBaseURL, cfg.LarkAccountsURL, cfg.LarkOAuthScope)
	repo := repository.NewAuthRepository(db)
	uc := usecase.NewAuthUsecase(repo, larkClient, cfg.RedirectURI)
	h := handler.NewAuthHandler(uc, cfg.ClientURL)

	return &AuthDomain{handler: h}
}

func (d *AuthDomain) RegisterRoutes(app *fiber.App) {
	authHTTP.RegisterAuthHTTPRoutes(app, d.handler)
}
