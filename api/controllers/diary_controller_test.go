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
	diary   *diary.Diary
	err     error
}

func (m *mockDiaryUsecase) FindAll() (*[]diary.Diary, error) {
	return m.diaries, m.err
}

func (m *mockDiaryUsecase) FindByUserID(userID int) (*[]diary.Diary, error) {
	return m.diaries, m.err
}

func (m *mockDiaryUsecase) FindByUserIDAndDate(userID int, date string) (*diary.Diary, error) {
	return m.diary, m.err
}

func (m *mockDiaryUsecase) FindByUserIDAndDateRange(userID int, startDate, endDate string) (*[]diary.Diary, error) {
	return m.diaries, m.err
}

func (m *mockDiaryUsecase) Create(diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryUsecase) Update(userID int, date string, diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryUsecase) Delete(userID int, date string) error {
	return m.err
}

// testDiaries はテスト用のダイアリーデータを定義します
var testDiaries = func() []diary.Diary {
	m5, _ := diary.NewMental(5)
	m3, _ := diary.NewMental(3)
	return []diary.Diary{
		{ID: 1, UserID: 1, Date: "2025-01-01", Mental: m5, Diary: "良い日だった"},
		{ID: 2, UserID: 1, Date: "2025-01-02", Mental: m3, Diary: "普通の日だった"},
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
		expectedData   []DiaryResponseDTO
		expectedError  string
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
			expectedData: []DiaryResponseDTO{
				{ID: 1, UserID: 1, Date: "2025-01-01", Mental: 5, Diary: "良い日だった"},
				{ID: 2, UserID: 1, Date: "2025-01-02", Mental: 3, Diary: "普通の日だった"},
			},
			expectedError: "",
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
			expectedData:   []DiaryResponseDTO{},
			expectedError:  "",
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
			expectedData:   nil,
			expectedError:  "DBエラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			controller := NewDiaryController(mock)

			router := gin.New()
			router.GET("/api/me/diaries", controller.FindAll)

			req, _ := http.NewRequest("GET", "/api/me/diaries", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response responseBody
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				// data部分を[]DiaryResponseDTOにデコードして比較
				dataBytes, _ := json.Marshal(response.Data)
				var actual []DiaryResponseDTO
				_ = json.Unmarshal(dataBytes, &actual)
				assert.Equal(t, tt.expectedData, actual)
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
		expectedData   *DiaryResponseDTO
		expectedError  string
	}{
		{
			name: "正常系：日記を作成できる",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			requestBody: CreateDiaryDTO{
				Date:   "2025-01-01",
				Mental: 5,
				Diary:  "良い日だった",
			},
			expectedStatus: http.StatusCreated,
			expectedData:   &DiaryResponseDTO{ID: 0, UserID: 1, Date: "2025-01-01", Mental: 5, Diary: "良い日だった"},
			expectedError:  "",
		},
		{
			name: "異常系：無効なメンタルスコア",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			requestBody: CreateDiaryDTO{
				Date:   "2025-01-01",
				Mental: 11, // 無効なメンタルスコア
				Diary:  "良い日だった",
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
			expectedError:  "mental value must be between 1 and 10",
		},
		{
			name: "異常系：複合ユニークキー制約違反",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("この日付の日記は既に作成されています"),
				}
			},
			requestBody: CreateDiaryDTO{
				Date:   "2025-01-01",
				Mental: 5,
				Diary:  "良い日だった",
			},
			expectedStatus: http.StatusConflict,
			expectedData:   nil,
			expectedError:  "この日付の日記は既に作成されています",
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("DBエラー"),
				}
			},
			requestBody: CreateDiaryDTO{
				Date:   "2025-01-01",
				Mental: 5,
				Diary:  "良い日だった",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
			expectedError:  "DBエラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			controller := NewDiaryController(mock)

			router := gin.New()
			router.POST("/api/me/diaries", controller.Create)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/me/diaries", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response responseBody
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				dataBytes, _ := json.Marshal(response.Data)
				var actual DiaryResponseDTO
				_ = json.Unmarshal(dataBytes, &actual)
				assert.Equal(t, *tt.expectedData, actual)
			}
		})
	}
}

func TestDiaryController_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func() *mockDiaryUsecase
		date           string
		requestBody    UpdateDiaryDTO
		expectedStatus int
		expectedData   *DiaryResponseDTO
		expectedError  string
	}{
		{
			name: "正常系：日記を更新できる",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			date: "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 7,
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusOK,
			expectedData:   &DiaryResponseDTO{ID: 0, UserID: 1, Date: "2025-01-01", Mental: 7, Diary: "更新された日記"},
			expectedError:  "",
		},
		{
			name: "異常系：無効なメンタルスコア",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			date: "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 11, // 無効なメンタルスコア
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
			expectedError:  "mental value must be between 1 and 10",
		},
		{
			name: "異常系：日記が見つからない",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("指定された日付の日記が見つかりません"),
				}
			},
			date: "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 7,
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusNotFound,
			expectedData:   nil,
			expectedError:  "指定された日付の日記が見つかりません",
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("DBエラー"),
				}
			},
			date: "2025-01-01",
			requestBody: UpdateDiaryDTO{
				Mental: 7,
				Diary:  "更新された日記",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
			expectedError:  "DBエラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			controller := NewDiaryController(mock)

			router := gin.New()
			router.PUT("/api/me/diaries/:date", controller.Update)

			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("PUT", "/api/me/diaries/"+tt.date, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response responseBody
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				dataBytes, _ := json.Marshal(response.Data)
				var actual DiaryResponseDTO
				_ = json.Unmarshal(dataBytes, &actual)
				assert.Equal(t, *tt.expectedData, actual)
			}
		})
	}
}

func TestDiaryController_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		setupMock       func() *mockDiaryUsecase
		date            string
		expectedStatus  int
		expectedMessage string
		expectedError   string
	}{
		{
			name: "正常系：日記を削除できる",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: nil,
				}
			},
			date:            "2025-01-01",
			expectedStatus:  http.StatusOK,
			expectedMessage: "日記が正常に削除されました",
			expectedError:   "",
		},
		{
			name: "異常系：日記が見つからない",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("指定された日付の日記が見つかりません"),
				}
			},
			date:            "2025-01-01",
			expectedStatus:  http.StatusNotFound,
			expectedMessage: "",
			expectedError:   "指定された日付の日記が見つかりません",
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("DBエラー"),
				}
			},
			date:            "2025-01-01",
			expectedStatus:  http.StatusInternalServerError,
			expectedMessage: "",
			expectedError:   "DBエラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			controller := NewDiaryController(mock)

			router := gin.New()
			router.DELETE("/api/me/diaries/:date", controller.Delete)

			req, _ := http.NewRequest("DELETE", "/api/me/diaries/"+tt.date, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response responseBody
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				dataBytes, _ := json.Marshal(response.Data)
				var data map[string]string
				_ = json.Unmarshal(dataBytes, &data)
				assert.Equal(t, tt.expectedMessage, data["message"])
			}
		})
	}
}

func TestDiaryController_FindByUserIDAndDate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func() *mockDiaryUsecase
		date           string
		expectedStatus int
		expectedData   *DiaryResponseDTO
		expectedError  string
	}{
		{
			name: "正常系：指定されたdateの日記を取得できる",
			setupMock: func() *mockDiaryUsecase {
				m5, _ := diary.NewMental(5)
				testDiary := diary.Diary{
					ID:     1,
					UserID: 1,
					Date:   "2025-01-01",
					Mental: m5,
					Diary:  "良い日だった",
				}
				return &mockDiaryUsecase{
					diary: &testDiary,
					err:   nil,
				}
			},
			date:           "2025-01-01",
			expectedStatus: http.StatusOK,
			expectedData:   &DiaryResponseDTO{ID: 1, UserID: 1, Date: "2025-01-01", Mental: 5, Diary: "良い日だった"},
			expectedError:  "",
		},
		{
			name: "異常系：日記が見つからない",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("指定された日付の日記が見つかりません"),
				}
			},
			date:           "2025-01-01",
			expectedStatus: http.StatusNotFound,
			expectedData:   nil,
			expectedError:  "指定された日付の日記が見つかりません",
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					err: errors.New("DBエラー"),
				}
			},
			date:           "2025-01-01",
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
			expectedError:  "DBエラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			controller := NewDiaryController(mock)

			router := gin.New()
			router.GET("/api/me/diaries/:date", controller.FindByUserIDAndDate)

			req, _ := http.NewRequest("GET", "/api/me/diaries/"+tt.date, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response responseBody
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				dataBytes, _ := json.Marshal(response.Data)
				var actual DiaryResponseDTO
				_ = json.Unmarshal(dataBytes, &actual)
				assert.Equal(t, *tt.expectedData, actual)
			}
		})
	}
}

func TestDiaryController_FindByUserIDAndDateRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func() *mockDiaryUsecase
		startDate      string
		endDate        string
		expectedStatus int
		expectedData   []DiaryResponseDTO
		expectedError  string
	}{
		{
			name: "正常系：指定された期間の日記を取得できる",
			setupMock: func() *mockDiaryUsecase {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryUsecase{
					diaries: &diaries,
					err:     nil,
				}
			},
			startDate:      "2025-01-01",
			endDate:        "2025-01-31",
			expectedStatus: http.StatusOK,
			expectedData: []DiaryResponseDTO{
				{ID: 1, UserID: 1, Date: "2025-01-01", Mental: 5, Diary: "良い日だった"},
				{ID: 2, UserID: 1, Date: "2025-01-02", Mental: 3, Diary: "普通の日だった"},
			},
			expectedError: "",
		},
		{
			name: "異常系：start_dateパラメータが不足",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					diaries: nil,
					err:     nil,
				}
			},
			startDate:      "",
			endDate:        "2025-01-31",
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
			expectedError:  "start_dateとend_dateの両方が必要です",
		},
		{
			name: "異常系：end_dateパラメータが不足",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					diaries: nil,
					err:     nil,
				}
			},
			startDate:      "2025-01-01",
			endDate:        "",
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
			expectedError:  "start_dateとend_dateの両方が必要です",
		},
		{
			name: "異常系：サービスがエラーを返した場合は500を返す",
			setupMock: func() *mockDiaryUsecase {
				return &mockDiaryUsecase{
					diaries: nil,
					err:     errors.New("DBエラー"),
				}
			},
			startDate:      "2025-01-01",
			endDate:        "2025-01-31",
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
			expectedError:  "DBエラー",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			controller := NewDiaryController(mock)

			router := gin.New()
			router.GET("/api/me/diaries/range", controller.FindByUserIDAndDateRange)

			url := "/api/me/diaries/range"
			if tt.startDate != "" {
				url += "?start_date=" + tt.startDate
			}
			if tt.endDate != "" {
				if tt.startDate != "" {
					url += "&end_date=" + tt.endDate
				} else {
					url += "?end_date=" + tt.endDate
				}
			}

			req, _ := http.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response responseBody
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError != "" {
				assert.Equal(t, tt.expectedError, response.Error)
			} else {
				dataBytes, _ := json.Marshal(response.Data)
				var actual []DiaryResponseDTO
				_ = json.Unmarshal(dataBytes, &actual)
				assert.Equal(t, tt.expectedData, actual)
			}
		})
	}
}
