package controllers

import (
	"emotra-backend/domain/diary"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// モックサービス
type mockDiaryUsecase struct {
	diaries *[]diary.Diary
	err     error
}

func (m *mockDiaryUsecase) FindAll() (*[]diary.Diary, error) {
	return m.diaries, m.err
}

// testDiaries はテスト用のダイアリーデータを定義します
var testDiaries = []diary.Diary{
	{ID: 1, UserID: 100, Date: "2025-01-01", Mental: diary.NewMental(5), Diary: "良い日だった"},
	{ID: 2, UserID: 101, Date: "2025-01-02", Mental: diary.NewMental(3), Diary: "普通の日だった"},
}

// responseBody はAPIレスポンスの構造を定義します
type responseBody struct {
	Data  []diary.Diary `json:"data"`
	Error string        `json:"error,omitempty"`
}

func TestDiaryController_FindAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func() *mockDiaryUsecase
		expectedStatus int
		expectedBody   responseBody
	}{
		{
			name: "正常系：全件取得できる場合は200を返す",
			setupMock: func() *mockDiaryUsecase {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryUsecase{
					diaries: &diaries,
					err:     nil,
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody: responseBody{
				Data: testDiaries,
			},
		},
		{
			name: "正常系：空のリストを処理できる",
			setupMock: func() *mockDiaryUsecase {
				emptyDiaries := make([]diary.Diary, 0)
				return &mockDiaryUsecase{
					diaries: &emptyDiaries,
					err:     nil,
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody: responseBody{
				Data: []diary.Diary{},
			},
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					diaries: nil,
					err:     errors.New("DBエラー"),
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: responseBody{
				Error: "DBエラー",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックのセットアップ
			usecase := tt.setupMock()
			controller := NewDiaryController(usecase)

			// テスト用のHTTPリクエストとレスポンス
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			// コントローラ呼び出し
			controller.FindAll(ctx)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code, "ステータスコードが期待値と異なります")

			// レスポンスボディの検証
			var response responseBody
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "レスポンスのJSONパースに失敗しました")

			if tt.expectedStatus == http.StatusOK {
				// 正常系の場合
				assert.Equal(t, tt.expectedBody.Data, response.Data, "レスポンスデータが期待値と異なります")
				assert.Empty(t, response.Error, "エラーメッセージが含まれています")
			} else {
				// 異常系の場合
				assert.Empty(t, response.Data, "データが含まれています")
				assert.Equal(t, tt.expectedBody.Error, response.Error, "エラーメッセージが期待値と異なります")
			}
		})
	}
}
