package controllers

import (
	"emotra-backend/api/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// モックサービス
type mockDiaryUsecase struct {
	diaries *[]models.Diary
	err     error
}

func (m *mockDiaryUsecase) FindAll() (*[]models.Diary, error) {
	return m.diaries, m.err
}

func TestDiaryController_FindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("正常に全件取得できる場合は200を返す", func(t *testing.T) {
		// モックデータ
		mockData := &[]models.Diary{
			{ID: 1, UserID: 100, Date: "2025-01-01", Mental: 4, Diary: "良い日だった"},
		}
		usecase := &mockDiaryUsecase{diaries: mockData}
		controller := NewDiaryController(usecase)

		// テスト用のHTTPリクエストとレスポンス
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		// コントローラ呼び出し
		controller.FindAll(ctx)

		// 検証
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), `"data"`)
		assert.Contains(t, w.Body.String(), `"良い日だった"`)
	})

	t.Run("サービスがエラーを返した場合は500を返す", func(t *testing.T) {
		usecase := &mockDiaryUsecase{err: errors.New("DBエラー")}
		controller := NewDiaryController(usecase)

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		controller.FindAll(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), `"error"`)
	})
}
