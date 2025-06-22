package main

import (
	"emotra-backend/infra"

	"emotra-backend/api/controllers"
	"emotra-backend/repositories"
	"emotra-backend/usecases"

	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	diaryRepository := repositories.NewDiaryRepository(db)
	diaryUsecase := usecases.NewDiaryUsecase(diaryRepository)
	diaryController := controllers.NewDiaryController(diaryUsecase)

	diaryAnalysisUsecase := usecases.NewDiaryAnalysisUsecase(diaryRepository)
	diaryAnalysisController := controllers.NewDiaryAnalysisController(diaryAnalysisUsecase)

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:3000" // デフォルト
	}

	// カンマ区切りで複数のオリジンを分割
	origins := strings.Split(corsOrigin, ",")
	// 空白を除去
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	router := gin.Default()

	// CORS設定を追加
	router.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	{
		api.GET("/me/diaries", diaryController.FindAll)
		api.GET("/me/diaries/:date", diaryController.FindByUserIDAndDate)
		api.POST("/me/diaries", diaryController.Create)
		api.PUT("/me/diaries/:date", diaryController.Update)
		api.DELETE("/me/diaries/:date", diaryController.Delete)
		api.GET("/me/analyze-diaries", diaryAnalysisController.AnalyzeAllDiariesHandler)
	}

	router.Run()
}
