package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	// รองรับหลาย origin (คั่นด้วย comma ใน env)
	allowOrigins := []string{frontendURL}
	if extra := os.Getenv("CORS_EXTRA_ORIGINS"); extra != "" {
		for _, o := range strings.Split(extra, ",") {
			allowOrigins = append(allowOrigins, strings.TrimSpace(o))
		}
	}

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
