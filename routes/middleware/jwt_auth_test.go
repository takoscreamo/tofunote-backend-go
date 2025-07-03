package middleware

import (
	"feelog-backend/domain/user"
	"feelog-backend/infra"
	"net/http"
	"net/http/httptest"
	"testing"

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
func (m *mockUserRepo) FindByEmail(email string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) FindByProviderId(provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Create(u *user.User) error { return nil }
func (m *mockUserRepo) Update(u *user.User) error { return nil }

func TestJWTAuthMiddleware_TableDriven(t *testing.T) {
	gin.SetMode(gin.TestMode)
	token, _ := infra.GenerateToken("00000000-0000-0000-0000-000000000042")

	tests := []struct {
		name       string
		header     string
		wantStatus int
		comment    string
	}{
		{"正常系: 有効なトークン", "Bearer " + token, 200, "有効なJWTで認証通過"},
		{"異常系: トークンなし", "", 401, "トークンなしで401"},
		{"異常系: 無効なトークン", "Bearer invalidtoken", 401, "不正なトークンで401"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(JWTAuthMiddleware())
			r.GET("/protected", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if !exists {
					c.JSON(500, gin.H{"error": "userID not set"})
					return
				}
				c.JSON(200, gin.H{"userID": userID})
			})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/protected", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code, tt.comment)
		})
	}
}

func TestJWTAuthMiddleware_GuestAndUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// ゲストユーザー
	guestID := uuid.New().String()
	guestToken, _ := infra.GenerateToken(guestID)
	SetAuthDB(nil)
	userRepo = &mockUserRepo{userByID: &user.User{ID: guestID, IsGuest: true}}

	r := gin.New()
	r.Use(JWTAuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		isGuest, _ := c.Get("isGuest")
		c.JSON(200, gin.H{"userID": userID, "isGuest": isGuest})
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+guestToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "true")

	// 正式ユーザー
	userID := uuid.New().String()
	userToken, _ := infra.GenerateToken(userID)
	userRepo = &mockUserRepo{userByID: &user.User{ID: userID, IsGuest: false}}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/protected", nil)
	req2.Header.Set("Authorization", "Bearer "+userToken)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
	assert.Contains(t, w2.Body.String(), "false")
}
