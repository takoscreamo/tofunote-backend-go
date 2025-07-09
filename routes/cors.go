package routes

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS CORS設定を設定
func SetupCORS(router *gin.Engine) {
	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		// 本番環境なら本番フロントのURLをデフォルトに
		if os.Getenv("ENV") == "prod" {
			corsOrigin = "https://feelog.takoscreamo.com"
		} else {
			corsOrigin = "http://localhost:3000"
		}
	}

	// カンマ区切りで複数のオリジンを分割
	origins := strings.Split(corsOrigin, ",")
	// 空白を除去
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	// CORS設定を追加
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}
