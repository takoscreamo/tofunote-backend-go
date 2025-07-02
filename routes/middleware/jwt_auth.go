package middleware

import (
	"net/http"
	"strings"

	"feelog-backend/infra"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証トークンが必要です"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		userID, err := infra.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}
		// userIDをコンテキストにセット
		c.Set("userID", userID)
		c.Next()
	}
}
