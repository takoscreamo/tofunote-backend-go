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

	router := gin.Default()
	router.GET("/diaries", diaryController.FindAll)
	router.Run()
}
