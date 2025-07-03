package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"feelog-backend/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	createdUser *user.User
	errCreate   error
}

func (m *mockUserRepo) FindByID(id string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) Create(u *user.User) error {
	if m.errCreate != nil {
		return m.errCreate
	}
	m.createdUser = u
	return nil
}
func (m *mockUserRepo) FindByProviderId(provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByRefreshToken(refreshToken string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) Update(u *user.User) error                                  { return nil }

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

			uc := NewUserController(tt.repo)
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
			} else {
				assert.Nil(t, tt.repo.createdUser)
			}
			assert.Contains(t, w.Body.String(), tt.wantBody)
		})
	}
}
