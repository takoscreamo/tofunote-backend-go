package db

import (
	"emotra-backend/domain/diary"

	"gorm.io/gorm"
)

type DiaryModel struct {
	gorm.Model
	UserID int    `gorm:"not null;type:integer;uniqueIndex:idx_user_date,priority:1"`
	Date   string `gorm:"not null;type:date;uniqueIndex:idx_user_date,priority:2"`
	Mental int    `gorm:"not null;type:integer"`
	Diary  string `gorm:"not null;type:text"`
}

func (DiaryModel) TableName() string {
	return "diaries"
}

// ToDomain converts the persistence model to the domain model.
func (d *DiaryModel) ToDomain() *diary.Diary {
	mental, err := diary.NewMental(d.Mental)
	if err != nil {
		// エラーの場合はデフォルト値を使用
		mental, _ = diary.NewMental(5)
	}
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
		Mental: d.Mental.GetValue(),
		Diary:  d.Diary,
	}
}
