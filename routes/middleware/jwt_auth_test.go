package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"tofunote-backend/domain/user"
	"tofunote-backend/infra"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	userByID *user.User
}

func (m *mockUserRepo) FindByID(id string) (*user.User, error) {
	return m.userByID, nil
}
func (m *mockUserRepo) FindByProviderId(provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Create(u *user.User) error                                  { return nil }
func (m *mockUserRepo) Update(u *user.User) error                                  { return nil }
func (m *mockUserRepo) FindByRefreshToken(refreshToken string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) DeleteByID(id string) error {
	return nil
}

func TestJWTAuthMiddleware_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testUUID := uuid.New().String()
	testToken, _ := infra.GenerateToken(testUUID)

	tests := []struct {
		name       string
		header     string
		mockUser   *user.User
		wantStatus int
		wantGuest  *bool // nilならisGuestチェックしない
		comment    string
	}{
		{
			name:       "有効なトークン(ゲスト)",
			header:     "Bearer " + testToken,
			mockUser:   &user.User{ID: testUUID, IsGuest: true},
			wantStatus: 200,
			wantGuest:  ptr(true),
			comment:    "ゲストユーザーで認証通過",
		},
		{
			name:       "有効なトークン(正式ユーザー)",
			header:     "Bearer " + testToken,
			mockUser:   &user.User{ID: testUUID, IsGuest: false},
			wantStatus: 200,
			wantGuest:  ptr(false),
			comment:    "正式ユーザーで認証通過",
		},
		{
			name:       "トークンなし",
			header:     "",
			mockUser:   nil,
			wantStatus: 401,
			wantGuest:  nil,
			comment:    "トークンなしで401",
		},
		{
			name:       "無効なトークン",
			header:     "Bearer invalidtoken",
			mockUser:   nil,
			wantStatus: 401,
			wantGuest:  nil,
			comment:    "不正なトークンで401",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAuthDB(nil)
			userRepo = &mockUserRepo{userByID: tt.mockUser}

			r := gin.New()
			r.Use(JWTAuthMiddleware())
			r.GET("/protected", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				isGuest, _ := c.Get("isGuest")
				c.JSON(200, gin.H{"userID": userID, "isGuest": isGuest})
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code, tt.comment)
			if tt.wantGuest != nil && w.Code == 200 {
				assert.Contains(t, w.Body.String(), "\"isGuest\":"+boolStr(*tt.wantGuest), tt.comment)
			}
		})
	}
}

func ptr(b bool) *bool { return &b }
func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
