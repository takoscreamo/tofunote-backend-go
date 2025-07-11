package db

import (
	"time"
)

type UserModel struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Nickname     string    `gorm:"type:varchar(255)"`
	Provider     string    `gorm:"type:varchar(50)"`
	ProviderID   string    `gorm:"type:varchar(255)"`
	IsGuest      bool      `gorm:"default:true"`
	RefreshToken string    `gorm:"type:varchar(255)"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (UserModel) TableName() string {
	return "users"
}
