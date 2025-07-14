package middleware

import (
	"net/http"
	"strings"

	"tofunote-backend/domain/user"
	"tofunote-backend/infra"
	"tofunote-backend/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB
var userRepo user.Repository

func SetAuthDB(db *gorm.DB) {
	dbInstance = db
	userRepo = repositories.NewUserRepository(db)
}

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
		// userIDからユーザー情報取得
		var isGuest bool
		if userRepo != nil {
			u, err := userRepo.FindByID(c.Request.Context(), userID)
			if err == nil && u != nil {
				isGuest = u.IsGuest
			}
		}
		c.Set("userID", userID)
		c.Set("isGuest", isGuest)
		c.Next()
	}
}
