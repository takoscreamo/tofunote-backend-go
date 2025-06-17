package main

import (
	"emotra-backend/infra"

	"emotra-backend/api/controllers"
	"emotra-backend/repositories"
	"emotra-backend/usecases"

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
	api := router.Group("/api")
	{
		api.GET("/diaries", diaryController.FindAll)
		api.POST("/diaries", diaryController.Create)
		api.GET("/analyze-diaries", diaryAnalysisController.AnalyzeAllDiariesHandler)
	}

	router.Run()
}
