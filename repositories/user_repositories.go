package repositories

import (
	"feelog-backend/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByProviderId(provider, providerId string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("provider = ? AND provider_id = ?", provider, providerId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByRefreshToken(refreshToken string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("refresh_token = ?", refreshToken).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(u *user.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepository) Create(u *user.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return r.db.Create(u).Error
}

var _ user.Repository = (*UserRepository)(nil)
