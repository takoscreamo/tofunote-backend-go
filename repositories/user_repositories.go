package repositories

import (
	"context"
	"tofunote-backend/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByProviderId(ctx context.Context, provider, providerId string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("provider = ? AND provider_id = ?", provider, providerId).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Where("refresh_token = ?", refreshToken).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *UserRepository) DeleteByID(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&user.User{}).Error
}
