package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	domainusecase "github.com/n2pluto/cinema-booking-system/internal/domain/usecase"
)

const oauthStateCookie = "oauth_state"

type AuthHandler struct {
	authUC domainusecase.AuthUseCase
}

func NewAuthHandler(authUC domainusecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

// GET /auth/google
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	// เก็บ state ใน httpOnly cookie เพื่อ verify ตอน callback
	c.SetCookie(oauthStateCookie, state, 300, "/", "", isSecure(), true)

	c.Redirect(http.StatusTemporaryRedirect, h.authUC.GetAuthURL(state))
}

// GET /auth/google/callback
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	// ตรวจ state ป้องกัน CSRF
	cookieState, err := c.Cookie(oauthStateCookie)
	if err != nil || cookieState == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state cookie"})
		return
	}

	if c.Query("state") != cookieState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state, possible CSRF attack"})
		return
	}

	// ลบ cookie ทันทีหลัง verify แล้ว
	c.SetCookie(oauthStateCookie, "", -1, "/", "", isSecure(), true)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	result, err := h.authUC.HandleCallback(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "authentication failed"})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/callback?token="+result.Token)
}

// GET /auth/me — ต้องแนบ JWT
func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetString("user_id"),
		"email":   c.GetString("email"),
		"role":    c.GetString("role"),
	})
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// isSecure returns true เมื่อ env = production (เพื่อ Secure cookie flag)
func isSecure() bool {
	return os.Getenv("APP_ENV") == "production"
}
