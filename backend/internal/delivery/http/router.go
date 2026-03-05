package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n2pluto/cinema-booking-system/internal/delivery/http/handler"
	"github.com/n2pluto/cinema-booking-system/internal/delivery/http/middleware"
	ws "github.com/n2pluto/cinema-booking-system/internal/delivery/http/ws"
	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
)

type Router struct {
	authHandler      *handler.AuthHandler
	cinemaHandler    *handler.CinemaHandler
	bookingHandler   *handler.BookingHandler
	auditLogHandler  *handler.AuditLogHandler
	websocketHandler *handler.WebSocketHandler
	hub              *ws.Hub
}

func NewRouter(
	authH *handler.AuthHandler,
	cinemaH *handler.CinemaHandler,
	bookingH *handler.BookingHandler,
	auditLogH *handler.AuditLogHandler,
	wsH *handler.WebSocketHandler,
	hub *ws.Hub,
) *Router {
	return &Router{
		authHandler:      authH,
		cinemaHandler:    cinemaH,
		bookingHandler:   bookingH,
		auditLogHandler:  auditLogH,
		websocketHandler: wsH,
		hub:              hub,
	}
}

func (ro *Router) Setup(r *gin.Engine) {
	// ─── Global middleware ────────────────────────────────────────────────────
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS())
	r.Use(middleware.RateLimit(20, 50))

	// ─── Health (public) ─────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// ─── OAuth (public) ──────────────────────────────────────────────────────
	r.GET("/auth/google", middleware.RateLimit(5, 10), ro.authHandler.GoogleLogin)
	r.GET("/auth/google/callback", middleware.RateLimit(5, 10), ro.authHandler.GoogleCallback)

	// ─── WebSocket (public route, auth via ?token query param) ───────────────
	r.GET("/api/ws/seats/:cinemaId", ro.websocketHandler.HandleSeatUpdates)

	// ─── Authenticated routes — USER และ ADMIN ────────────────────────────────
	api := r.Group("/api")
	api.Use(middleware.JWTAuth())
	api.Use(middleware.RequireRole(entity.RoleUser, entity.RoleAdmin))
	{
		// Auth
		api.GET("/me", ro.authHandler.Me)
		api.GET("/auth/me", ro.authHandler.Me)

		// Cinema (read-only — ทั้ง USER และ ADMIN ดูได้)
		api.GET("/cinema", ro.cinemaHandler.List)
		api.GET("/seats/:cinemaId", ro.cinemaHandler.GetSeats)
	}

	// ─── Booking routes — เฉพาะ USER เท่านั้น (admin ไม่มีสิทธิ์จอง) ─────────
	userOnly := r.Group("/api")
	userOnly.Use(middleware.JWTAuth())
	userOnly.Use(middleware.RequireRole(entity.RoleUser))
	{
		userOnly.POST("/booking/lock", ro.bookingHandler.Lock)
		userOnly.POST("/booking/confirm", ro.bookingHandler.Confirm)
		userOnly.POST("/booking/cancel", ro.bookingHandler.Cancel)
		userOnly.GET("/booking/me", ro.bookingHandler.GetMine)
	}

	// ─── Admin routes — เฉพาะ ADMIN ──────────────────────────────────────────
	admin := r.Group("/api/admin")
	admin.Use(middleware.JWTAuth())
	admin.Use(middleware.RequireRole(entity.RoleAdmin))
	{
		admin.GET("/cinema", ro.cinemaHandler.List)
		admin.POST("/cinema", ro.cinemaHandler.CreateShowtime)
		admin.GET("/audit-logs", ro.auditLogHandler.List)
		admin.GET("/bookings", ro.bookingHandler.ListAll)
	}
}
