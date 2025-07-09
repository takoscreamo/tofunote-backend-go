package routes

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupCORS CORS設定を設定
func SetupCORS(router *gin.Engine) {
	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000" // デフォルト
	}

	// カンマ区切りで複数のオリジンを分割
	origins := strings.Split(corsOrigin, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""
		for _, o := range origins {
			if o == origin {
				allowOrigin = origin
				break
			}
		}
		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			c.Writer.Header().Set("Vary", "Origin")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}
