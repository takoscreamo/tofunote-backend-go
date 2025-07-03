package user

import "time"

// Userエンティティ
type User struct {
	ID           string
	Provider     string
	ProviderID   string
	IsGuest      bool
	RefreshToken string
	CreatedAt    time.Time
}
