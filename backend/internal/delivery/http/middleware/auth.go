package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/n2pluto/cinema-booking-system/internal/domain/entity"
	jwtpkg "github.com/n2pluto/cinema-booking-system/pkg/jwt"
)

// JWTAuth validates the Bearer token and injects claims into context.
// ทุก route ที่ใช้ middleware นี้ต้องแนบ header:
//
//	Authorization: Bearer <token>
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization format, use: Bearer <token>"})
			return
		}

		claims, err := jwtpkg.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("display_name", claims.DisplayName)
		c.Next()
	}
}

// RequireRole อนุญาตเฉพาะ role ที่ระบุ (ใช้ต่อจาก JWTAuth เสมอ)
func RequireRole(roles ...entity.UserRole) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[string(r)] = true
	}

	return func(c *gin.Context) {
		role := c.GetString("role")
		if !allowed[role] {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden: insufficient role"})
			return
		}
		c.Next()
	}
}
