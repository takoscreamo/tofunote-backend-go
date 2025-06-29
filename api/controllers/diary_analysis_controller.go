package controllers

import (
	"net/http"

	"feelog-backend/usecases"

	"github.com/gin-gonic/gin"
)

type DiaryAnalysisController struct {
	DiaryAnalysisUsecase *usecases.DiaryAnalysisUsecase
}

// NewDiaryAnalysisController は新しい DiaryAnalysisController を作成する
func NewDiaryAnalysisController(usecase *usecases.DiaryAnalysisUsecase) *DiaryAnalysisController {
	return &DiaryAnalysisController{
		DiaryAnalysisUsecase: usecase,
	}
}

// AnalyzeAllDiariesHandler は認証されたユーザーの日記を分析するエンドポイント
func (c *DiaryAnalysisController) AnalyzeAllDiariesHandler(ctx *gin.Context) {
	// ハードコードでuser_id=1を使用
	userID := 1
	result, err := c.DiaryAnalysisUsecase.AnalyzeUserDiaries(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"analysis_result": result})
}
