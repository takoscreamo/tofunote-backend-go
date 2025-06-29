package main

import (
	"feelog-backend/infra"
	"feelog-backend/routes"

	"feelog-backend/api/controllers"
	"feelog-backend/repositories"
	"feelog-backend/usecases"

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

	router := gin.Default()

	// ヘルスチェックエンドポイントを追加
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// CORS設定を追加
	routes.SetupCORS(router)

	// SwaggerUIとOpenAPI仕様書のエンドポイントを設定
	routes.SetupSwaggerEndpoints(router)

	// APIエンドポイントを設定
	routes.SetupAPIEndpoints(router, diaryController, diaryAnalysisController)

	router.Run()
}
