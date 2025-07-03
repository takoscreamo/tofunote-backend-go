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
)

type mockUserRepo struct {
	createdUser *user.User
	findByID    *user.User
}

func (m *mockUserRepo) FindByID(id string) (*user.User, error) {
	return m.findByID, nil
}
func (m *mockUserRepo) Create(u *user.User) error {
	m.createdUser = u
	return nil
}
func (m *mockUserRepo) FindByProviderId(provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Update(u *user.User) error { return nil }

func TestGuestLogin_NewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockUserRepo{}
	uc := NewUserController(repo)
	r := gin.New()
	r.POST("/api/guest-login", uc.GuestLogin)

	uuid := "b3b1a2e0-4b5c-4e2a-8c2e-1b2c3d4e5f60"
	body := map[string]string{"id": uuid}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/guest-login", bytes.NewReader(b))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.NotNil(t, repo.createdUser)
	assert.Equal(t, uuid, repo.createdUser.ID)
}

func TestGuestLogin_ExistingUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	uuid := "b3b1a2e0-4b5c-4e2a-8c2e-1b2c3d4e5f60"
	existing := &user.User{ID: uuid, IsGuest: true}
	repo := &mockUserRepo{findByID: existing}
	uc := NewUserController(repo)
	r := gin.New()
	r.POST("/api/guest-login", uc.GuestLogin)

	body := map[string]string{"id": existing.ID}
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/guest-login", bytes.NewReader(b))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())
	assert.Nil(t, repo.createdUser) // 新規作成されない
}
