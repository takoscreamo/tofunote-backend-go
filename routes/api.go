package routes

import (
	"feelog-backend/api/controllers"
	"log"
	"os"
	"time"

	"feelog-backend/routes/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// SetupAPIEndpoints APIエンドポイントを設定
func SetupAPIEndpoints(router *gin.Engine, diaryController *controllers.DiaryController, diaryAnalysisController *controllers.DiaryAnalysisController, userController *controllers.UserController) {
	// ヘルスチェックエンドポイント
	router.GET("/ping", func(c *gin.Context) {
		log.Printf("[DEBUG] Ping endpoint called - returning pong message")
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 詳細な動作確認エンドポイント
	router.GET("/health", func(c *gin.Context) {
		log.Printf("[DEBUG] Health endpoint called - returning detailed health info")

		// 環境変数の確認（機密情報は除外）
		env := os.Getenv("ENV")
		if env == "" {
			env = "dev"
		}

		healthInfo := gin.H{
			"status":      "healthy",
			"timestamp":   time.Now().Format(time.RFC3339),
			"environment": env,
			"service":     "feelog-backend",
			"version":     "1.0.0",
			"endpoint":    "/health",
			"message":     "This is the health check endpoint",
			"headers":     c.Request.Header,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
		}

		c.JSON(200, healthInfo)
	})

	// シンプルな動作確認エンドポイント
	router.GET("/status", func(c *gin.Context) {
		log.Printf("[DEBUG] Status endpoint called - returning simple status")
		c.JSON(200, gin.H{
			"status":    "ok",
			"message":   "Service is running",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	api := router.Group("/api")
	{
		api.POST("/guest-login", userController.GuestLogin)
		api.POST("/refresh-token", userController.RefreshToken)

		// 認証が必要なグループ
		auth := api.Group("")
		auth.Use(middleware.JWTAuthMiddleware())
		auth.GET("/me/diaries", diaryController.FindAll)
		auth.GET("/me/diaries/range", diaryController.FindByUserIDAndDateRange)
		auth.GET("/me/diaries/:date", diaryController.FindByUserIDAndDate)
		auth.POST("/me/diaries", diaryController.Create)
		auth.PUT("/me/diaries/:date", diaryController.Update)
		auth.DELETE("/me/diaries/:date", diaryController.Delete)
		auth.GET("/me/analyze-diaries", diaryAnalysisController.AnalyzeAllDiariesHandler)
		auth.DELETE("/me", userController.DeleteMe)
		auth.GET("/me", userController.GetMe)
		auth.PATCH("/me", userController.PatchMe)
	}
}
