package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"feelog-backend/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	userByEmail *user.User
	err         error
	createdUser *user.User
}

func (m *mockUserRepo) FindByEmail(email string) (*user.User, error) {
	return m.userByEmail, m.err
}
func (m *mockUserRepo) FindByProviderId(provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Create(u *user.User) error {
	m.createdUser = u
	return m.err
}

func TestUserController_Register_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		registerBody   RegisterRequest
		registerRepo   *mockUserRepo
		registerStatus int
		loginBody      LoginRequest
		loginRepo      *mockUserRepo
		loginStatus    int
		comment        string
	}{
		{
			name:           "正常系: 登録→ログイン",
			registerBody:   RegisterRequest{Email: "test@example.com", Password: "password", Nickname: "nick"},
			registerRepo:   &mockUserRepo{},
			registerStatus: http.StatusOK,
			loginBody:      LoginRequest{Email: "test@example.com", Password: "password"},
			loginRepo: func() *mockUserRepo {
				hash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
				return &mockUserRepo{userByEmail: &user.User{ID: 1, Email: "test@example.com", PasswordHash: string(hash)}}
			}(),
			loginStatus: http.StatusOK,
			comment:     "正常な登録とログインができる",
		},
		{
			name:           "異常系: 登録時に重複",
			registerBody:   RegisterRequest{Email: "dup@example.com", Password: "password", Nickname: "nick"},
			registerRepo:   &mockUserRepo{userByEmail: &user.User{Email: "dup@example.com"}},
			registerStatus: http.StatusConflict,
			comment:        "既存メールで登録すると409",
		},
		{
			name:        "異常系: ログイン失敗（パスワード不一致）",
			loginBody:   LoginRequest{Email: "test@example.com", Password: "wrongpass"},
			loginRepo:   &mockUserRepo{userByEmail: &user.User{Email: "test@example.com", PasswordHash: "$2a$10$invalidhash"}},
			loginStatus: http.StatusUnauthorized,
			comment:     "パスワード不一致で401",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 登録テスト
			if tt.registerRepo != nil {
				r := gin.New()
				uc := NewUserController(tt.registerRepo)
				r.POST("/register", uc.Register)
				b, _ := json.Marshal(tt.registerBody)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/register", bytes.NewReader(b))
				r.ServeHTTP(w, req)
				assert.Equal(t, tt.registerStatus, w.Code, tt.comment)
			}
			// ログインテスト
			if tt.loginRepo != nil {
				r := gin.New()
				uc := NewUserController(tt.loginRepo)
				r.POST("/login", uc.Login)
				b, _ := json.Marshal(tt.loginBody)
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("POST", "/login", bytes.NewReader(b))
				r.ServeHTTP(w, req)
				assert.Equal(t, tt.loginStatus, w.Code, tt.comment)
			}
		})
	}
}
