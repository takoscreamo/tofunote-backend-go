package models

import "gorm.io/gorm"

type Diary struct {
	gorm.Model
	UserID int    `gorm:"not null;type:integer"`
	Date   string `gorm:"not null;type:date"`
	Mental int    `gorm:"not null;type:integer"`
	Diary  string `gorm:"not null;type:text"`
}
