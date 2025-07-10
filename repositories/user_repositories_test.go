package repositories

import (
	"testing"
	"tofunote-backend/domain/user"

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

func TestUserRepository_Create(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)
	tests := []struct {
		name     string
		toCreate *user.User
	}{
		{"Create guest user", &user.User{ID: "uuid-1", IsGuest: true}},
		{"Create oauth user", &user.User{ID: "uuid-2", Provider: "google", ProviderID: "gid-2", IsGuest: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(tt.toCreate)
			assert.NoError(t, err)
			found, err := repo.FindByID(tt.toCreate.ID)
			assert.NoError(t, err)
			assert.NotNil(t, found)
			assert.Equal(t, tt.toCreate.ID, found.ID)
			assert.Equal(t, tt.toCreate.IsGuest, found.IsGuest)
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)
	repo.Create(&user.User{ID: "uuid-1", IsGuest: true})
	repo.Create(&user.User{ID: "uuid-2", Provider: "google", ProviderID: "gid-2", IsGuest: false})

	tests := []struct {
		name     string
		findID   string
		expectID *string
	}{
		{"Found guest user", "uuid-1", ptrStr("uuid-1")},
		{"Found oauth user", "uuid-2", ptrStr("uuid-2")},
		{"Not found", "not-exist", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByID(tt.findID)
			assert.NoError(t, err)
			if tt.expectID != nil {
				assert.NotNil(t, found)
				assert.Equal(t, *tt.expectID, found.ID)
			} else {
				assert.Nil(t, found)
			}
		})
	}
}

func TestUserRepository_FindByProviderId(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)
	repo.Create(&user.User{ID: "uuid-1", Provider: "google", ProviderID: "gid-1", IsGuest: false})
	repo.Create(&user.User{ID: "uuid-2", Provider: "apple", ProviderID: "aid-2", IsGuest: false})

	tests := []struct {
		name       string
		provider   string
		providerID string
		expectID   *string
	}{
		{"Found google user", "google", "gid-1", ptrStr("uuid-1")},
		{"Found apple user", "apple", "aid-2", ptrStr("uuid-2")},
		{"Not found", "google", "not-exist", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindByProviderId(tt.provider, tt.providerID)
			assert.NoError(t, err)
			if tt.expectID != nil {
				assert.NotNil(t, found)
				assert.Equal(t, *tt.expectID, found.ID)
			} else {
				assert.Nil(t, found)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupUserTestDB(t)
	repo := NewUserRepository(db)
	repo.Create(&user.User{ID: "uuid-1", IsGuest: true})

	t.Run("Update IsGuest", func(t *testing.T) {
		u, _ := repo.FindByID("uuid-1")
		u.IsGuest = false
		err := repo.Update(u)
		assert.NoError(t, err)
		updated, _ := repo.FindByID("uuid-1")
		assert.Equal(t, false, updated.IsGuest)
	})
}

func ptrStr(s string) *string { return &s }
