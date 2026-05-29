package handler

import (
	"errors"
	"log/slog"
	"net/url"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"ufriend-cx-dashboard-server/internal/auth/usecase"
	"ufriend-cx-dashboard-server/pkg/response"
)

const (
	sessionCookieName = "ufriend_session"
	oauthStateCookie  = "lark_oauth_state"
)

type AuthHandler struct {
	usecase   usecase.AuthUsecase
	clientURL string
	isSecure  bool
}

func NewAuthHandler(uc usecase.AuthUsecase, clientURL string) *AuthHandler {
	isSecure := os.Getenv("APP_ENV") == "production"
	return &AuthHandler{
		usecase:   uc,
		clientURL: clientURL,
		isSecure:  isSecure,
	}
}

// LarkRedirect ส่งผู้ใช้ไปหน้า authorize ของ Lark
func (h *AuthHandler) LarkRedirect(c *fiber.Ctx) error {
	state := uuid.NewString()
	c.Cookie(&fiber.Cookie{
		Name:     oauthStateCookie,
		Value:    state,
		HTTPOnly: true,
		Secure:   h.isSecure,
		SameSite: "Lax",
		MaxAge:   600,
		Path:     "/",
	})

	url := h.usecase.GetAuthorizeURL(state)
	return c.Redirect(url, fiber.StatusFound)
}

// LarkCallback รับ code จาก Lark แลก token แล้ว redirect กลับ client
func (h *AuthHandler) LarkCallback(c *fiber.Ctx) error {
	// ผู้ใช้ปฏิเสธการ authorize — Lark ส่ง error=access_denied
	if oauthErr := c.Query("error"); oauthErr != "" {
		if oauthErr == "access_denied" {
			return h.redirectWithError(c, "คุณปฏิเสธการเข้าสู่ระบบ")
		}
		return h.redirectWithError(c, "เข้าสู่ระบบไม่สำเร็จ")
	}

	code := c.Query("code")
	if code == "" {
		return h.redirectWithError(c, "ไม่พบ authorization code")
	}

	// ตรวจสอบ state เพื่อป้องกัน CSRF
	state := c.Query("state")
	storedState := c.Cookies(oauthStateCookie)
	if state == "" || storedState == "" || state != storedState {
		return h.redirectWithError(c, "state ไม่ถูกต้อง")
	}
	c.Cookie(&fiber.Cookie{Name: oauthStateCookie, Value: "", MaxAge: -1, Path: "/"})

	sessionID, err := h.usecase.HandleCallback(c.Context(), code)
	if err != nil {
		slog.Error("lark callback failed", "error", err)
		return h.redirectWithError(c, "เข้าสู่ระบบไม่สำเร็จ")
	}

	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		HTTPOnly: true,
		Secure:   h.isSecure,
		SameSite: "Lax",
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
		Path:     "/",
	})

	return c.Redirect(h.clientURL+"/login?login=success", fiber.StatusFound)
}

// GetMe ดึง profile ผู้ใช้จาก session ปัจจุบัน
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	sessionID := c.Cookies(sessionCookieName)
	if sessionID == "" {
		return response.Error(c, fiber.StatusUnauthorized, "ไม่พบเซสชัน", "ERR_UNAUTHORIZED")
	}

	me, err := h.usecase.GetMe(c.Context(), sessionID)
	if err != nil {
		if errors.Is(err, usecase.ErrSessionNotFound) || errors.Is(err, usecase.ErrSessionExpired) {
			h.clearSessionCookie(c)
			return response.Error(c, fiber.StatusUnauthorized, "เซสชันหมดอายุ กรุณาเข้าสู่ระบบใหม่", "ERR_SESSION_EXPIRED")
		}
		slog.Error("get me failed", "error", err)
		return response.ErrorServer(c, "ไม่สามารถดึงข้อมูลผู้ใช้ได้", err)
	}

	return response.OK(c, "ดึงข้อมูลผู้ใช้สำเร็จ", me)
}

// Logout ลบ session และ cookie
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionID := c.Cookies(sessionCookieName)
	if sessionID != "" {
		if err := h.usecase.Logout(c.Context(), sessionID); err != nil {
			slog.Error("logout failed", "error", err)
		}
	}
	h.clearSessionCookie(c)
	return response.OK(c, "ออกจากระบบสำเร็จ", nil)
}

func (h *AuthHandler) clearSessionCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		HTTPOnly: true,
		Secure:   h.isSecure,
		SameSite: "Lax",
		MaxAge:   -1,
		Path:     "/",
	})
}

func (h *AuthHandler) redirectWithError(c *fiber.Ctx, msg string) error {
	return c.Redirect(h.clientURL+"/login?error="+url.QueryEscape(msg), fiber.StatusFound)
}
