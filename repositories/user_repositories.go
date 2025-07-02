package repositories

import (
	"feelog-backend/domain/user"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	var u user.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
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

func (r *UserRepository) Create(u *user.User) error {
	return r.db.Create(u).Error
}

var _ user.Repository = (*UserRepository)(nil)
