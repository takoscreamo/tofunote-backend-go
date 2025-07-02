package user

import "time"

// Userエンティティ
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	Nickname     string
	Provider     string
	ProviderID   string
	CreatedAt    time.Time
}
