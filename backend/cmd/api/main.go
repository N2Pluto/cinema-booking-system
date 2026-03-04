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
	mongorepo "github.com/n2pluto/cinema-booking-system/internal/repository/mongodb"
	"github.com/n2pluto/cinema-booking-system/internal/usecase/auth"
	"github.com/n2pluto/cinema-booking-system/pkg/database"
)

func main() {
	// ─── Infrastructure ───────────────────────────────────────────────────────
	mongodb, err := database.NewMongoDB()
	if err != nil {
		log.Fatalf("MongoDB: %v", err)
	}
	defer mongodb.Disconnect()

	// ─── Repositories ─────────────────────────────────────────────────────────
	userRepo := mongorepo.NewUserRepository(mongodb.Collection("users"))

	// ─── Use Cases ────────────────────────────────────────────────────────────
	authUC := auth.NewUseCase(
		userRepo,
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_CALLBACK_URL"),
	)

	// ─── Handlers ─────────────────────────────────────────────────────────────
	authH := handler.NewAuthHandler(authUC)

	// ─── Router ───────────────────────────────────────────────────────────────
	r := gin.Default()
	delivery.NewRouter(authH).Setup(r)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
