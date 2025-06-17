package controllers

import (
	"emotra-backend/domain/diary"
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": diaries})
}

type CreateDiaryDTO struct {
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Mental int    `json:"mental"`
	Diary  string `json:"diary"`
}

func (c *DiaryController) Create(ctx *gin.Context) {
	var req CreateDiaryDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストデータです"})
		return
	}
	mental, err := diary.NewMental(req.Mental)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newDiary := diary.Diary{
		UserID: req.UserID,
		Date:   req.Date,
		Mental: mental,
		Diary:  req.Diary,
	}
	err = c.usecase.Create(&newDiary)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": newDiary})
}
