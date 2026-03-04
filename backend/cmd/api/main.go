package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	delivery "github.com/n2pluto/cinema-booking-system/internal/delivery/http"
	"github.com/n2pluto/cinema-booking-system/internal/delivery/http/handler"
	wsHub "github.com/n2pluto/cinema-booking-system/internal/delivery/http/ws"
	mongorepo "github.com/n2pluto/cinema-booking-system/internal/repository/mongodb"
	auditlogusecase "github.com/n2pluto/cinema-booking-system/internal/usecase/audit_log"
	"github.com/n2pluto/cinema-booking-system/internal/usecase/auth"
	bookingusecase "github.com/n2pluto/cinema-booking-system/internal/usecase/booking"
	cinemausecase "github.com/n2pluto/cinema-booking-system/internal/usecase/cinema"
	"github.com/n2pluto/cinema-booking-system/pkg/database"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ─── Infrastructure ───────────────────────────────────────────────────────
	mongodb, err := database.NewMongoDB()
	if err != nil {
		log.Fatalf("MongoDB: %v", err)
	}
	defer mongodb.Disconnect()

	// ─── WebSocket Hub ────────────────────────────────────────────────────────
	hub := wsHub.NewHub()
	go hub.Run()

	// ─── Repositories ─────────────────────────────────────────────────────────
	userRepo := mongorepo.NewUserRepository(mongodb.Collection("users"))
	cinemaRepo := mongorepo.NewCinemaRepository(mongodb.Collection("cinemas"))
	bookingRepo := mongorepo.NewBookingRepository(mongodb.Collection("bookings"))
	auditRepo := mongorepo.NewAuditLogRepository(mongodb.Collection("audit_logs"))

	// ─── Use Cases ────────────────────────────────────────────────────────────
	authUC := auth.NewUseCase(
		userRepo,
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_CALLBACK_URL"),
	)
	cinemaUC := cinemausecase.NewUseCase(cinemaRepo)
	bookingUC := bookingusecase.NewUseCase(cinemaRepo, bookingRepo, auditRepo, hub)
	auditUC := auditlogusecase.NewUseCase(auditRepo)

	// Start background seat-lock timeout ticker
	bookingUC.StartTimeoutTicker(ctx)

	// ─── Handlers ─────────────────────────────────────────────────────────────
	authH := handler.NewAuthHandler(authUC)
	cinemaH := handler.NewCinemaHandler(cinemaUC)
	bookingH := handler.NewBookingHandler(bookingUC)
	auditH := handler.NewAuditLogHandler(auditUC)
	wsH := handler.NewWebSocketHandler(hub)

	// ─── Router ───────────────────────────────────────────────────────────────
	r := gin.Default()
	delivery.NewRouter(authH, cinemaH, bookingH, auditH, wsH, hub).Setup(r)

	// ─── Server ───────────────────────────────────────────────────────────────
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		log.Printf("Server running on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(shutCtx)
}
