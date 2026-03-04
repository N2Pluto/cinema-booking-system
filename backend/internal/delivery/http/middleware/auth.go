package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
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
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid authorization format"})
			return
		}

		claims, err := jwtpkg.Parse(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid or expired token"})
			return
		}

		// inject claims เข้า context ให้ handler ใช้ได้
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("user", claims)

		c.Next()
	}
}

// AdminOnly ใช้ต่อจาก JWTAuth — อนุญาตเฉพาะ role ADMIN
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "ADMIN" {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
