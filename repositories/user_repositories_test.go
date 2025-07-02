package repositories

import (
	"feelog-backend/domain/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	if err := db.AutoMigrate(&user.User{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	return db
}

func TestUserRepository_Create_FindByEmail(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)
	tests := []struct {
		name        string
		toCreate    *user.User
		findEmail   string
		expectFound bool
	}{
		{"create and find", &user.User{Email: "a@a.com", PasswordHash: "hash", Nickname: "nick"}, "a@a.com", true},
		{"not found", nil, "notfound@example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.toCreate != nil {
				err := repo.Create(tt.toCreate)
				assert.NoError(t, err)
				assert.NotZero(t, tt.toCreate.ID)
			}
			found, err := repo.FindByEmail(tt.findEmail)
			assert.NoError(t, err)
			if tt.expectFound {
				assert.NotNil(t, found)
				assert.Equal(t, tt.findEmail, found.Email)
			} else {
				assert.Nil(t, found)
			}
		})
	}
}

func TestUserRepository_FindByProviderId(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)
	tests := []struct {
		name        string
		toCreate    *user.User
		provider    string
		providerId  string
		expectFound bool
	}{
		{"create and find", &user.User{Email: "b@b.com", PasswordHash: "hash", Nickname: "nick", Provider: "google", ProviderID: "gid"}, "google", "gid", true},
		{"not found", nil, "google", "notexist", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.toCreate != nil {
				repo.Create(tt.toCreate)
			}
			found, err := repo.FindByProviderId(tt.provider, tt.providerId)
			assert.NoError(t, err)
			if tt.expectFound {
				assert.NotNil(t, found)
				assert.Equal(t, tt.toCreate.Email, found.Email)
			} else {
				assert.Nil(t, found)
			}
		})
	}
}
