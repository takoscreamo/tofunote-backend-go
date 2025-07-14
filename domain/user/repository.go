package user

import "context"

// UserRepositoryはユーザーの永続化インターフェース
//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=user

type Repository interface {
	FindByProviderId(ctx context.Context, provider, providerId string) (*User, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*User, error)
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	DeleteByID(ctx context.Context, id string) error
}
