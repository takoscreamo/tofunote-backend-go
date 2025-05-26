package main

import (
	"emotra-backend/api/controllers"
	"emotra-backend/api/models"
	"emotra-backend/api/repositories"
	"emotra-backend/api/usecases"

	"github.com/gin-gonic/gin"
)

func main() {
	diaries := []models.Diary{
		{ID: 1, UserID: 1, Date: "2025-05-01", Mental: 5, Diary: "今日はとても良い日でした。"},
		{ID: 2, UserID: 1, Date: "2025-05-02", Mental: 3, Diary: "普通の日でした。"},
		{ID: 3, UserID: 2, Date: "2025-05-01", Mental: 4, Diary: "少し疲れましたが、充実した日でした。"},
	}

	diaryRepository := repositories.NewDiaryMemoryRepository(diaries)
	diaryUsecase := usecases.NewDiaryUsecase(diaryRepository)
	diaryController := controllers.NewDiaryController(diaryUsecase)

	router := gin.Default()
	router.GET("/diaries", diaryController.FindAll)
	router.Run()
}
