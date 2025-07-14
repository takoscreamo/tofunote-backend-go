package controllers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tofunote-backend/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	createdUser  *user.User
	errCreate    error
	FindByIDFunc func(ctx context.Context, id string) (*user.User, error)
	UpdateFunc   func(ctx context.Context, u *user.User) error
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*user.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	if m.errCreate != nil {
		return m.errCreate
	}
	m.createdUser = u
	return nil
}
func (m *mockUserRepo) FindByProviderId(ctx context.Context, provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByRefreshToken(ctx context.Context, refreshToken string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, u)
	}
	return nil
}
func (m *mockUserRepo) DeleteByID(ctx context.Context, id string) error {
	return nil
}

func (m *mockUserRepo) setUser(u *user.User) {
	m.createdUser = u
}

func TestGuestLogin_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		repo         *mockUserRepo
		mockTokenErr bool
		wantStatus   int
		wantUser     bool
		wantBody     string
	}{
		{
			name:       "正常系: 新規ゲスト作成",
			repo:       &mockUserRepo{},
			wantStatus: http.StatusOK,
			wantUser:   true,
			wantBody:   "token",
		},
		{
			name:       "異常系: Create失敗",
			repo:       &mockUserRepo{errCreate: errors.New("fail")},
			wantStatus: http.StatusInternalServerError,
			wantUser:   false,
			wantBody:   "ユーザー作成に失敗しました",
		},
		{
			name:         "異常系: トークン生成失敗",
			repo:         &mockUserRepo{},
			mockTokenErr: true,
			wantStatus:   http.StatusInternalServerError,
			wantUser:     true,
			wantBody:     "トークン生成に失敗しました",
		},
	}

	origGenerateToken := generateTokenForTest
	defer func() { generateTokenForTest = origGenerateToken }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generateTokenForTest = func(id string) (string, error) {
				if tt.mockTokenErr {
					return "", errors.New("fail")
				}
				return "dummy-token", nil
			}

			uc := NewUserController(tt.repo, nil) // withdrawUsecaseは不要なためnilでOK
			r := gin.New()
			r.POST("/api/guest-login", func(c *gin.Context) {
				uc.GuestLogin(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/guest-login", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code, w.Body.String())
			if tt.wantUser {
				assert.NotNil(t, tt.repo.createdUser)
				assert.NotEmpty(t, tt.repo.createdUser.ID)
				assert.True(t, tt.repo.createdUser.IsGuest)
				assert.Equal(t, "ゲスト", tt.repo.createdUser.Nickname)
			} else {
				assert.Nil(t, tt.repo.createdUser)
			}
			assert.Contains(t, w.Body.String(), tt.wantBody)
		})
	}
}

func TestGetMe_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		findByIDFunc func(ctx context.Context, id string) (*user.User, error)
	}

	tests := []struct {
		name       string
		fields     fields
		userID     interface{}
		wantStatus int
		wantBody   string
	}{
		{
			name: "正常系: ユーザー情報取得",
			fields: fields{
				findByIDFunc: func(ctx context.Context, id string) (*user.User, error) {
					return &user.User{ID: id, Nickname: "テスト太郎"}, nil
				},
			},
			userID:     "test-id",
			wantStatus: http.StatusOK,
			wantBody:   "テスト太郎",
		},
		{
			name:       "異常系: 認証情報なし",
			fields:     fields{},
			userID:     nil,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "認証情報が見つかりません",
		},
		{
			name: "異常系: ユーザー未発見",
			fields: fields{
				findByIDFunc: func(ctx context.Context, id string) (*user.User, error) {
					return nil, nil
				},
			},
			userID:     "notfound",
			wantStatus: http.StatusInternalServerError,
			wantBody:   "ユーザーが見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepo{FindByIDFunc: tt.fields.findByIDFunc}
			uc := NewUserController(mockRepo, nil)
			r := gin.New()
			r.GET("/api/me", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("userID", tt.userID)
				}
				uc.GetMe(c)
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/me", nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantBody)
		})
	}
}

func TestPatchMe_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type fields struct {
		findByIDFunc func(ctx context.Context, id string) (*user.User, error)
		updateFunc   func(ctx context.Context, u *user.User) error
	}

	tests := []struct {
		name       string
		fields     fields
		userID     interface{}
		body       string
		wantStatus int
		wantBody   string
		wantNick   string
	}{
		{
			name: "正常系: ニックネーム更新",
			fields: fields{
				findByIDFunc: func(ctx context.Context, id string) (*user.User, error) {
					return &user.User{ID: id, Nickname: "旧名"}, nil
				},
				updateFunc: func(ctx context.Context, u *user.User) error {
					return nil
				},
			},
			userID:     "test-id",
			body:       `{"nickname": "新しい名"}`,
			wantStatus: http.StatusOK,
			wantBody:   "新しい名",
			wantNick:   "新しい名",
		},
		{
			name:       "異常系: 認証情報なし",
			fields:     fields{},
			userID:     nil,
			body:       `{"nickname": "新しい名"}`,
			wantStatus: http.StatusUnauthorized,
			wantBody:   "認証情報が見つかりません",
		},
		{
			name: "異常系: ユーザー未発見",
			fields: fields{
				findByIDFunc: func(ctx context.Context, id string) (*user.User, error) {
					return nil, nil
				},
			},
			userID:     "notfound",
			body:       `{"nickname": "新しい名"}`,
			wantStatus: http.StatusInternalServerError,
			wantBody:   "ユーザーが見つかりません",
		},
		{
			name: "異常系: バリデーションエラー（空ボディ）",
			fields: fields{
				findByIDFunc: func(ctx context.Context, id string) (*user.User, error) {
					return &user.User{ID: id, Nickname: "旧名"}, nil
				},
			},
			userID:     "test-id",
			body:       ``,
			wantStatus: http.StatusBadRequest,
			wantBody:   "無効なリクエストデータです",
		},
		{
			name: "異常系: 更新可能な項目なし",
			fields: fields{
				findByIDFunc: func(ctx context.Context, id string) (*user.User, error) {
					return &user.User{ID: id, Nickname: "旧名"}, nil
				},
			},
			userID:     "test-id",
			body:       `{"not_exist": "xxx"}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   "更新可能な項目がありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockUserRepo{
				FindByIDFunc: tt.fields.findByIDFunc,
				UpdateFunc:   tt.fields.updateFunc,
			}
			uc := NewUserController(mockRepo, nil)
			r := gin.New()
			r.PATCH("/api/me", func(c *gin.Context) {
				if tt.userID != nil {
					c.Set("userID", tt.userID)
				}
				uc.PatchMe(c)
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PATCH", "/api/me", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Contains(t, w.Body.String(), tt.wantBody)
			if tt.wantNick != "" && tt.wantStatus == http.StatusOK {
				// レスポンスのnicknameが正しいか確認
				assert.Contains(t, w.Body.String(), tt.wantNick)
			}
		})
	}
}
