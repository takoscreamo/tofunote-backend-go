package user

// UserRepositoryはユーザーの永続化インターフェース
//go:generate mockgen -source=repository.go -destination=mock_repository.go -package=user

type Repository interface {
	FindByProviderId(provider, providerId string) (*User, error)
	FindByRefreshToken(refreshToken string) (*User, error)
	Create(user *User) error
	FindByID(id string) (*User, error)
	Update(user *User) error
}
