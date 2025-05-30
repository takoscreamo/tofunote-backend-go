package db

import (
	"emotra-backend/domain/diary"

	"gorm.io/gorm"
)

type DiaryModel struct {
	gorm.Model
	UserID int    `gorm:"not null;type:integer"`
	Date   string `gorm:"not null;type:date"`
	Mental int    `gorm:"not null;type:integer"`
	Diary  string `gorm:"not null;type:text"`
}

func (DiaryModel) TableName() string {
	return "diaries"
}

// ToDomain converts the persistence model to the domain model.
func (d *DiaryModel) ToDomain() *diary.Diary {
	return &diary.Diary{
		ID:     int(d.ID),
		UserID: d.UserID,
		Date:   d.Date,
		Mental: diary.Mental{Value: d.Mental},
		Diary:  d.Diary,
	}
}

// FromDomain converts the domain model to the persistence model.
func FromDomain(d *diary.Diary) *DiaryModel {
	return &DiaryModel{
		UserID: d.UserID,
		Date:   d.Date,
		Mental: d.Mental.Value,
		Diary:  d.Diary,
	}
}
