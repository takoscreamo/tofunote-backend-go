package controllers

import (
	"emotra-backend/domain/diary"
	"emotra-backend/usecases"
	"net/http"
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
	// ハードコードでuser_id=1を使用
	userID := 1
	diaries, err := c.usecase.FindByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responseDTOs := make([]DiaryResponseDTO, 0, len(*diaries))
	for _, d := range *diaries {
		responseDTOs = append(responseDTOs, ToResponseDTO(&d))
	}

	ctx.JSON(http.StatusOK, gin.H{"data": responseDTOs})
}

func (c *DiaryController) FindByUserIDAndDate(ctx *gin.Context) {
	// ハードコードでuser_id=1を使用
	userID := 1
	date := ctx.Param("date")

	diary, err := c.usecase.FindByUserIDAndDate(userID, date)
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
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Date   string `json:"date"`
	Mental int    `json:"mental"`
	Diary  string `json:"diary"`
}

// ToResponseDTO converts domain Diary to response DTO
func ToResponseDTO(diary *diary.Diary) DiaryResponseDTO {
	return DiaryResponseDTO{
		ID:     diary.ID,
		UserID: diary.UserID,
		Date:   diary.Date,
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

	// ハードコードでuser_id=1を使用
	userID := 1
	newDiary := diary.Diary{
		UserID: userID,
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
	ctx.JSON(http.StatusCreated, gin.H{"data": ToResponseDTO(&newDiary)})
}

func (c *DiaryController) Update(ctx *gin.Context) {
	// ハードコードでuser_id=1を使用
	userID := 1
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

	ctx.JSON(http.StatusOK, gin.H{"data": ToResponseDTO(&updateDiary)})
}

func (c *DiaryController) Delete(ctx *gin.Context) {
	// ハードコードでuser_id=1を使用
	userID := 1
	date := ctx.Param("date")

	err := c.usecase.Delete(userID, date)
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
