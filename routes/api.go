package routes

import (
	"feelog-backend/api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupAPIEndpoints APIエンドポイントを設定
func SetupAPIEndpoints(router *gin.Engine, diaryController *controllers.DiaryController, diaryAnalysisController *controllers.DiaryAnalysisController) {
	// ヘルスチェックエンドポイント
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := router.Group("/api")
	{
		api.GET("/me/diaries", diaryController.FindAll)
		api.GET("/me/diaries/range", diaryController.FindByUserIDAndDateRange)
		api.GET("/me/diaries/:date", diaryController.FindByUserIDAndDate)
		api.POST("/me/diaries", diaryController.Create)
		api.PUT("/me/diaries/:date", diaryController.Update)
		api.DELETE("/me/diaries/:date", diaryController.Delete)
		api.GET("/me/analyze-diaries", diaryAnalysisController.AnalyzeAllDiariesHandler)
	}
}
