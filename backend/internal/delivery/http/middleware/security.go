package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders adds HTTP security headers to every response.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ป้องกัน Clickjacking
		c.Header("X-Frame-Options", "DENY")
		// ป้องกัน MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		// ป้องกัน XSS บน browser เก่า
		c.Header("X-XSS-Protection", "1; mode=block")
		// บังคับ HTTPS (ปิดไว้ก่อนถ้า dev ไม่มี TLS)
		// c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		// ลด info ที่ส่งใน Referer header
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// จำกัด browser features ที่ไม่จำเป็น
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		// ซ่อน server info
		c.Header("Server", "")

		c.Next()
	}
}
