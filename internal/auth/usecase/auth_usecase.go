package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/auth/model"
)

const refreshThreshold = 5 * time.Minute

func (u *authUsecase) GetAuthorizeURL(state string) string {
	return u.larkClient.AuthorizeURL(u.redirectURI, state)
}

// HandleCallback แลก code เป็น token แล้วบันทึก session ใหม่
func (u *authUsecase) HandleCallback(ctx context.Context, code string) (string, error) {
	appToken, err := u.larkClient.GetAppAccessToken()
	if err != nil {
		return "", err
	}

	userToken, err := u.larkClient.ExchangeCode(appToken, code)
	if err != nil {
		return "", err
	}

	sessionID := uuid.NewString()
	now := time.Now()
	session := &model.AuthSession{
		SessionID:    sessionID,
		UserID:       userToken.UserID,
		AccessToken:  userToken.AccessToken,
		RefreshToken: userToken.RefreshToken,
		ExpiresAt:    now.Add(time.Duration(userToken.ExpiresIn) * time.Second),
		CreatedAt:    now,
	}

	if err := u.repo.Save(ctx, session); err != nil {
		return "", err
	}

	return sessionID, nil
}

// GetMe ดึง profile — auto-refresh ถ้า token ใกล้หมดอายุ (< 5 นาที)
func (u *authUsecase) GetMe(ctx context.Context, sessionID string) (*UserMeResponse, error) {
	session, err := u.repo.FindBySessionID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	session, err = u.ensureFreshToken(ctx, session)
	if err != nil {
		return nil, err
	}

	profile, err := u.larkClient.GetUserInfo(session.AccessToken)
	if err != nil {
		return nil, err
	}

	userID := profile.UserID
	if userID == "" {
		userID = profile.OpenID
	}

	return &UserMeResponse{
		Name:      profile.Name,
		EnName:    profile.EnName,
		AvatarURL: profile.AvatarURL,
		OpenID:    profile.OpenID,
		Email:     profile.Email,
		UserID:    userID,
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, sessionID string) error {
	return u.repo.DeleteBySessionID(ctx, sessionID)
}

// ensureFreshToken ตรวจสอบและ refresh token ถ้าใกล้หมดอายุ
func (u *authUsecase) ensureFreshToken(ctx context.Context, session *model.AuthSession) (*model.AuthSession, error) {
	if time.Until(session.ExpiresAt) > refreshThreshold {
		return session, nil
	}

	appToken, err := u.larkClient.GetAppAccessToken()
	if err != nil {
		return nil, err
	}

	refreshed, err := u.larkClient.RefreshAccessToken(appToken, session.RefreshToken)
	if err != nil {
		return nil, ErrSessionExpired
	}

	session.AccessToken = refreshed.AccessToken
	session.RefreshToken = refreshed.RefreshToken
	session.ExpiresAt = time.Now().Add(time.Duration(refreshed.ExpiresIn) * time.Second)

	if err := u.repo.Save(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Sentinel errors สำหรับ handler map เป็น HTTP status
var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
)
