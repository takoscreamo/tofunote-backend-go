package controllers

import (
	"emotra-backend/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DiaryController struct {
	usecase usecases.IDiaryUsecase
}

func NewDiaryController(usecase usecases.IDiaryUsecase) *DiaryController {
	return &DiaryController{usecase: usecase}
}

func (c *DiaryController) FindAll(ctx *gin.Context) {
	diaries, err := c.usecase.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": diaries})
}
