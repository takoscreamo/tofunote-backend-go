package controllers

import (
	"emotra-backend/domain/diary"
	"emotra-backend/usecases"
	"net/http"
	"strconv"
	"strings"

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

type UpdateDiaryDTO struct {
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
		// 複合ユニークキー制約違反の場合は409 Conflictを返す
		if strings.Contains(err.Error(), "この日付の日記は既に作成されています") {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": newDiary})
}

func (c *DiaryController) Update(ctx *gin.Context) {
	// URLパラメータからuser_idとdateを取得
	userIDStr := ctx.Param("user_id")
	date := ctx.Param("date")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーIDです"})
		return
	}

	var req UpdateDiaryDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なリクエストデータです"})
		return
	}

	mental, err := diary.NewMental(req.Mental)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateDiary := diary.Diary{
		UserID: userID,
		Date:   date,
		Mental: mental,
		Diary:  req.Diary,
	}

	err = c.usecase.Update(userID, date, &updateDiary)
	if err != nil {
		if strings.Contains(err.Error(), "指定された日付の日記が見つかりません") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updateDiary})
}

func (c *DiaryController) Delete(ctx *gin.Context) {
	// URLパラメータからuser_idとdateを取得
	userIDStr := ctx.Param("user_id")
	date := ctx.Param("date")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なユーザーIDです"})
		return
	}

	err = c.usecase.Delete(userID, date)
	if err != nil {
		if strings.Contains(err.Error(), "指定された日付の日記が見つかりません") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "日記が正常に削除されました"})
}
