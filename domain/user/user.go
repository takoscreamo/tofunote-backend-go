package user

import "time"

// Userエンティティ
type User struct {
	ID           string
	Email        string
	PasswordHash string
	Nickname     string
	Provider     string
	ProviderID   string
	IsGuest      bool
	CreatedAt    time.Time
}
