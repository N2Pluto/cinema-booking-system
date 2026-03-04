package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n2pluto/cinema-booking-system/internal/delivery/http/handler"
	"github.com/n2pluto/cinema-booking-system/internal/delivery/http/middleware"
)

type Router struct {
	authHandler *handler.AuthHandler
}

func NewRouter(authHandler *handler.AuthHandler) *Router {
	return &Router{authHandler: authHandler}
}

func (ro *Router) Setup(r *gin.Engine) {
	// ─── Global middleware ────────────────────────────────────────────────────
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(20, 50)) // 20 req/s per IP, burst 50

	// ─── Health ──────────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ─── Auth (public) ───────────────────────────────────────────────────────
	authGroup := r.Group("/auth")
	authGroup.Use(middleware.RateLimit(5, 10)) // auth เข้มงวดกว่า: 5 req/s per IP
	{
		authGroup.GET("/google", ro.authHandler.GoogleLogin)
		authGroup.GET("/google/callback", ro.authHandler.GoogleCallback)
	}

	// ─── Protected — ต้องแนบ JWT ทุก request ─────────────────────────────────
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/auth/me", ro.authHandler.Me)
	}
}
