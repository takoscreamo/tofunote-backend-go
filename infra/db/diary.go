package db

import (
	"emotra-backend/domain/diary"

	"gorm.io/gorm"
)

type DiaryModel struct {
	gorm.Model
	UserID int    `gorm:"not null;type:integer;uniqueIndex:idx_user_date,priority:1" json:"user_id"`
	Date   string `gorm:"not null;type:date;uniqueIndex:idx_user_date,priority:2" json:"date"`
	Mental int    `gorm:"not null;type:integer" json:"mental"`
	Diary  string `gorm:"not null;type:text" json:"diary"`
}

func (DiaryModel) TableName() string {
	return "diaries"
}

// ToDomain converts the persistence model to the domain model.
func (d *DiaryModel) ToDomain() *diary.Diary {
	mental, _ := diary.NewMental(d.Mental)
	return &diary.Diary{
		ID:     int(d.ID),
		UserID: d.UserID,
		Date:   d.Date,
		Mental: mental,
		Diary:  d.Diary,
	}
}

// FromDomain converts the domain model to the persistence model.
func FromDomain(d *diary.Diary) *DiaryModel {
	return &DiaryModel{
		UserID: d.UserID,
		Date:   d.Date,
		Mental: int(d.Mental),
		Diary:  d.Diary,
	}
}
