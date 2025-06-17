package controllers

import (
	"bytes"
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

func (m *mockDiaryUsecase) Create(diary *diary.Diary) error {
	return m.err
}

// testDiaries はテスト用のダイアリーデータを定義します
var testDiaries = func() []diary.Diary {
	m5, _ := diary.NewMental(5)
	m3, _ := diary.NewMental(3)
	return []diary.Diary{
		{ID: 1, UserID: 100, Date: "2025-01-01", Mental: m5, Diary: "良い日だった"},
		{ID: 2, UserID: 101, Date: "2025-01-02", Mental: m3, Diary: "普通の日だった"},
	}
}()

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

func TestDiaryController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func() *mockDiaryUsecase
		requestBody    CreateDiaryDTO
		expectedStatus int
		expectedBody   responseBody
	}{
		{
			name: "正常系：日記を作成できる",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			requestBody: CreateDiaryDTO{
				UserID: 100,
				Date:   "2025-01-01",
				Mental: 5,
				Diary:  "良い日だった",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: responseBody{
				Data: []diary.Diary{
					{
						UserID: 100,
						Date:   "2025-01-01",
						Mental: diary.Mental{Value: 5},
						Diary:  "良い日だった",
					},
				},
			},
		},
		{
			name: "異常系：無効なリクエストデータ",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			requestBody: CreateDiaryDTO{
				UserID: 100,
				Date:   "2025-01-01",
				Mental: 11, // 無効なメンタルスコア
				Diary:  "良い日だった",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: responseBody{
				Error: "mental value must be between 1 and 10",
			},
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("DBエラー"),
				}
			},
			requestBody: CreateDiaryDTO{
				UserID: 100,
				Date:   "2025-01-01",
				Mental: 5,
				Diary:  "良い日だった",
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

			// リクエストボディの設定
			jsonData, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}
			ctx.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
			ctx.Request.Header.Set("Content-Type", "application/json")

			// コントローラ呼び出し
			controller.Create(ctx)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code, "ステータスコードが期待値と異なります")

			// レスポンスボディの検証
			var response map[string]json.RawMessage
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "レスポンスのJSONパースに失敗しました")

			if tt.expectedStatus == http.StatusCreated {
				// 正常系の場合
				var diary diary.Diary
				_ = json.Unmarshal(response["data"], &diary)
				assert.Equal(t, tt.expectedBody.Data[0].UserID, diary.UserID, "ユーザーIDが期待値と異なります")
				assert.Equal(t, tt.expectedBody.Data[0].Date, diary.Date, "日付が期待値と異なります")
				assert.Equal(t, tt.expectedBody.Data[0].Mental.Value, diary.Mental.Value, "メンタルスコアが期待値と異なります")
				assert.Equal(t, tt.expectedBody.Data[0].Diary, diary.Diary, "日記内容が期待値と異なります")
			} else {
				// 異常系の場合
				var errResp map[string]string
				_ = json.Unmarshal(w.Body.Bytes(), &errResp)
				assert.Equal(t, tt.expectedBody.Error, errResp["error"], "エラーメッセージが期待値と異なります")
			}
		})
	}
}
