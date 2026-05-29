package usecase

import (
	"context"

	"ufriend-cx-dashboard-server/internal/auth/lark"
	"ufriend-cx-dashboard-server/internal/auth/repository"
)

// UserMeResponse ข้อมูลผู้ใช้ที่ส่งกลับ client (camelCase แปลงที่ handler/dto)
type UserMeResponse struct {
	Name      string `json:"name"`
	EnName    string `json:"en_name"`
	AvatarURL string `json:"avatar_url"`
	OpenID    string `json:"open_id"`
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
}

type AuthUsecase interface {
	GetAuthorizeURL(state string) string
	HandleCallback(ctx context.Context, code string) (sessionID string, err error)
	GetMe(ctx context.Context, sessionID string) (*UserMeResponse, error)
	Logout(ctx context.Context, sessionID string) error
}

type authUsecase struct {
	repo        repository.AuthRepository
	larkClient  *lark.Client
	redirectURI string
}

func NewAuthUsecase(repo repository.AuthRepository, larkClient *lark.Client, redirectURI string) AuthUsecase {
	return &authUsecase{
		repo:        repo,
		larkClient:  larkClient,
		redirectURI: redirectURI,
	}
}
