package controllers

import (
	"net/http"
	"strings"
	"time"
	"tofunote-backend/domain/diary"
	"tofunote-backend/usecases"

	"github.com/gin-gonic/gin"
)

type DiaryController struct {
	usecase usecases.IDiaryUsecase
}

func NewDiaryController(usecase usecases.IDiaryUsecase) *DiaryController {
	return &DiaryController{usecase: usecase}
}

func (c *DiaryController) FindAll(ctx *gin.Context) {
	// JWTトークンからuserIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報が見つかりません"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの形式が不正です"})
		return
	}
	diaries, err := c.usecase.FindByUserID(ctx.Request.Context(), userIDStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseDTOs := make([]DiaryResponseDTO, 0, len(diaries))
	for _, d := range diaries {
		responseDTOs = append(responseDTOs, ToResponseDTO(&d))
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseDTOs})
}

func (c *DiaryController) FindByUserIDAndDate(ctx *gin.Context) {
	// JWTトークンからuserIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報が見つかりません"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの形式が不正です"})
		return
	}
	date := ctx.Param("date")

	diary, err := c.usecase.FindByUserIDAndDate(ctx.Request.Context(), userIDStr, date)
	if err != nil {
		if strings.Contains(err.Error(), "指定された日付の日記が見つかりません") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": ToResponseDTO(diary)})
}

func (c *DiaryController) FindByUserIDAndDateRange(ctx *gin.Context) {
	// JWTトークンからuserIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報が見つかりません"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの形式が不正です"})
		return
	}
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")

	// クエリパラメータのバリデーション
	if startDate == "" || endDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_dateとend_dateの両方が必要です"})
		return
	}

	diaries, err := c.usecase.FindByUserIDAndDateRange(ctx.Request.Context(), userIDStr, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseDTOs := make([]DiaryResponseDTO, 0, len(diaries))
	for _, d := range diaries {
		responseDTOs = append(responseDTOs, ToResponseDTO(&d))
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseDTOs})
}

type CreateDiaryDTO struct {
	Date   string `json:"date"`
	Mental int    `json:"mental"`
	Diary  string `json:"diary"`
}

type UpdateDiaryDTO struct {
	Mental int    `json:"mental"`
	Diary  string `json:"diary"`
}

type DiaryResponseDTO struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Date   string `json:"date"`
	Mental int    `json:"mental"`
	Diary  string `json:"diary"`
}

// ToResponseDTO converts domain Diary to response DTO
func ToResponseDTO(diary *diary.Diary) DiaryResponseDTO {
	date := diary.Date
	if t, err := time.Parse(time.RFC3339, diary.Date); err == nil {
		date = t.Format("2006-01-02")
	}
	return DiaryResponseDTO{
		ID:     diary.ID,
		UserID: diary.UserID,
		Date:   date,
		Mental: int(diary.Mental),
		Diary:  diary.Diary,
	}
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

	// JWTトークンからuserIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報が見つかりません"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの形式が不正です"})
		return
	}
	newDiary := diary.Diary{
		UserID: userIDStr,
		Date:   req.Date,
		Mental: mental,
		Diary:  req.Diary,
	}
	err = c.usecase.Create(ctx.Request.Context(), &newDiary)
	if err != nil {
		// 複合ユニークキー制約違反の場合は409 Conflictを返す
		if strings.Contains(err.Error(), "この日付の日記は既に作成されています") {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": ToResponseDTO(&newDiary)})
}

func (c *DiaryController) Update(ctx *gin.Context) {
	// JWTトークンからuserIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報が見つかりません"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの形式が不正です"})
		return
	}
	date := ctx.Param("date")

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
		UserID: userIDStr,
		Date:   date,
		Mental: mental,
		Diary:  req.Diary,
	}

	err = c.usecase.Update(ctx.Request.Context(), userIDStr, date, &updateDiary)
	if err != nil {
		if strings.Contains(err.Error(), "指定された日付の日記が見つかりません") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": ToResponseDTO(&updateDiary)})
}

func (c *DiaryController) Delete(ctx *gin.Context) {
	// JWTトークンからuserIDを取得
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証情報が見つかりません"})
		return
	}
	userIDStr, ok := userID.(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーIDの形式が不正です"})
		return
	}
	date := ctx.Param("date")

	err := c.usecase.Delete(ctx.Request.Context(), userIDStr, date)
	if err != nil {
		if strings.Contains(err.Error(), "指定された日付の日記が見つかりません") {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "日記が正常に削除されました"}})
}
