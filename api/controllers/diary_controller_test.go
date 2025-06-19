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

func (m *mockDiaryUsecase) Update(userID int, date string, diary *diary.Diary) error {
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
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
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
			var response map[string]json.RawMessage
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "レスポンスのJSONパースに失敗しました")

			if tt.expectedStatus == http.StatusOK {
				// 正常系の場合
				var actualDiaries []diary.Diary
				err = json.Unmarshal(response["data"], &actualDiaries)
				assert.NoError(t, err, "Diary配列のJSONパースに失敗しました")
				assert.Equal(t, len(tt.expectedBody.Data.([]diary.Diary)), len(actualDiaries), "配列の長さが期待値と異なります")
				for i, expected := range tt.expectedBody.Data.([]diary.Diary) {
					if i < len(actualDiaries) {
						assert.Equal(t, expected.UserID, actualDiaries[i].UserID, "ユーザーIDが期待値と異なります")
						assert.Equal(t, expected.Date, actualDiaries[i].Date, "日付が期待値と異なります")
						assert.Equal(t, expected.Mental.GetValue(), actualDiaries[i].Mental.GetValue(), "メンタルスコアが期待値と異なります")
						assert.Equal(t, expected.Diary, actualDiaries[i].Diary, "日記内容が期待値と異なります")
					}
				}
			} else {
				// 異常系の場合
				var errorResponse responseBody
				err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "エラーレスポンスのJSONパースに失敗しました")
				assert.Equal(t, tt.expectedBody.Error, errorResponse.Error, "エラーメッセージが期待値と異なります")
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
				Data: diary.Diary{
					UserID: 100,
					Date:   "2025-01-01",
					Mental: func() diary.Mental { m, _ := diary.NewMental(5); return m }(),
					Diary:  "良い日だった",
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
			name: "異常系：複合ユニークキー制約違反",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("この日付の日記は既に作成されています"),
				}
			},
			requestBody: CreateDiaryDTO{
				UserID: 100,
				Date:   "2025-01-01",
				Mental: 5,
				Diary:  "良い日だった",
			},
			expectedStatus: http.StatusConflict,
			expectedBody: responseBody{
				Error: "この日付の日記は既に作成されています",
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

			if tt.expectedBody.Error != "" {
				var errorResponse responseBody
				err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "エラーレスポンスのJSONパースに失敗しました")
				assert.Equal(t, tt.expectedBody.Error, errorResponse.Error, "エラーメッセージが期待値と異なります")
			} else {
				// 成功レスポンスの場合は、期待値の型に応じて比較
				if expectedDiary, ok := tt.expectedBody.Data.(diary.Diary); ok {
					// 単一のDiaryオブジェクトの場合
					var actualDiary diary.Diary
					err = json.Unmarshal(response["data"], &actualDiary)
					assert.NoError(t, err, "DiaryオブジェクトのJSONパースに失敗しました")
					assert.Equal(t, expectedDiary.UserID, actualDiary.UserID, "ユーザーIDが期待値と異なります")
					assert.Equal(t, expectedDiary.Date, actualDiary.Date, "日付が期待値と異なります")
					assert.Equal(t, expectedDiary.Mental.GetValue(), actualDiary.Mental.GetValue(), "メンタルスコアが期待値と異なります")
					assert.Equal(t, expectedDiary.Diary, actualDiary.Diary, "日記内容が期待値と異なります")
				} else if expectedDiaries, ok := tt.expectedBody.Data.([]diary.Diary); ok {
					// Diary配列の場合
					var actualDiaries []diary.Diary
					err = json.Unmarshal(response["data"], &actualDiaries)
					assert.NoError(t, err, "Diary配列のJSONパースに失敗しました")
					assert.Equal(t, len(expectedDiaries), len(actualDiaries), "配列の長さが期待値と異なります")
					for i, expected := range expectedDiaries {
						if i < len(actualDiaries) {
							assert.Equal(t, expected.UserID, actualDiaries[i].UserID, "ユーザーIDが期待値と異なります")
							assert.Equal(t, expected.Date, actualDiaries[i].Date, "日付が期待値と異なります")
							assert.Equal(t, expected.Mental.GetValue(), actualDiaries[i].Mental.GetValue(), "メンタルスコアが期待値と異なります")
							assert.Equal(t, expected.Diary, actualDiaries[i].Diary, "日記内容が期待値と異なります")
						}
					}
				}
			}
		})
	}
}

func TestDiaryController_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func() *mockDiaryUsecase
		userID         string
		date           string
		requestBody    UpdateDiaryDTO
		expectedStatus int
		expectedBody   responseBody
	}{
		{
			name: "正常系：日記を更新できる",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			userID: "100",
			date:   "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 7,
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusOK,
			expectedBody: responseBody{
				Data: diary.Diary{
					UserID: 100,
					Date:   "2025-01-01",
					Mental: func() diary.Mental { m, _ := diary.NewMental(7); return m }(),
					Diary:  "更新された日記",
				},
			},
		},
		{
			name: "異常系：無効なユーザーID",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			userID: "invalid",
			date:   "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 7,
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: responseBody{
				Error: "無効なユーザーIDです",
			},
		},
		{
			name: "異常系：無効なメンタルスコア",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			userID: "100",
			date:   "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 11, // 無効なメンタルスコア
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: responseBody{
				Error: "mental value must be between 1 and 10",
			},
		},
		{
			name: "異常系：日記が見つからない",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("指定された日付の日記が見つかりません"),
				}
			},
			userID: "100",
			date:   "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 7,
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: responseBody{
				Error: "指定された日付の日記が見つかりません",
			},
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("DBエラー"),
				}
			},
			userID: "100",
			date:   "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 7,
				Diary:  "更新された日記",
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

			// パラメータの設定
			ctx.Params = gin.Params{
				{Key: "user_id", Value: tt.userID},
				{Key: "date", Value: tt.date},
			}

			// リクエストボディの設定
			jsonData, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}
			ctx.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer(jsonData))
			ctx.Request.Header.Set("Content-Type", "application/json")

			// コントローラ呼び出し
			controller.Update(ctx)

			// ステータスコードの検証
			assert.Equal(t, tt.expectedStatus, w.Code, "ステータスコードが期待値と異なります")

			// レスポンスボディの検証
			var response map[string]json.RawMessage
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err, "レスポンスのJSONパースに失敗しました")

			if tt.expectedBody.Error != "" {
				var errorResponse responseBody
				err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err, "エラーレスポンスのJSONパースに失敗しました")
				assert.Equal(t, tt.expectedBody.Error, errorResponse.Error, "エラーメッセージが期待値と異なります")
			} else {
				// 成功レスポンスの場合は、期待値の型に応じて比較
				if expectedDiary, ok := tt.expectedBody.Data.(diary.Diary); ok {
					// 単一のDiaryオブジェクトの場合
					var actualDiary diary.Diary
					err = json.Unmarshal(response["data"], &actualDiary)
					assert.NoError(t, err, "DiaryオブジェクトのJSONパースに失敗しました")
					assert.Equal(t, expectedDiary.UserID, actualDiary.UserID, "ユーザーIDが期待値と異なります")
					assert.Equal(t, expectedDiary.Date, actualDiary.Date, "日付が期待値と異なります")
					assert.Equal(t, expectedDiary.Mental.GetValue(), actualDiary.Mental.GetValue(), "メンタルスコアが期待値と異なります")
					assert.Equal(t, expectedDiary.Diary, actualDiary.Diary, "日記内容が期待値と異なります")
				} else if expectedDiaries, ok := tt.expectedBody.Data.([]diary.Diary); ok {
					// Diary配列の場合
					var actualDiaries []diary.Diary
					err = json.Unmarshal(response["data"], &actualDiaries)
					assert.NoError(t, err, "Diary配列のJSONパースに失敗しました")
					assert.Equal(t, len(expectedDiaries), len(actualDiaries), "配列の長さが期待値と異なります")
					for i, expected := range expectedDiaries {
						if i < len(actualDiaries) {
							assert.Equal(t, expected.UserID, actualDiaries[i].UserID, "ユーザーIDが期待値と異なります")
							assert.Equal(t, expected.Date, actualDiaries[i].Date, "日付が期待値と異なります")
							assert.Equal(t, expected.Mental.GetValue(), actualDiaries[i].Mental.GetValue(), "メンタルスコアが期待値と異なります")
							assert.Equal(t, expected.Diary, actualDiaries[i].Diary, "日記内容が期待値と異なります")
						}
					}
				}
			}
		})
	}
}
