package db

import (
	"database/sql/driver"
	"emotra-backend/domain/diary"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// DateString is a custom type for handling PostgreSQL date type
type DateString string

// Scan implements the sql.Scanner interface
func (d *DateString) Scan(value interface{}) error {
	if value == nil {
		*d = ""
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*d = DateString(v.Format("2006-01-02"))
	case string:
		*d = DateString(v)
	case []byte:
		*d = DateString(string(v))
	default:
		return fmt.Errorf("cannot scan %T into DateString", value)
	}
	return nil
}

// Value implements the driver.Valuer interface
func (d DateString) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}
	return string(d), nil
}

type DiaryModel struct {
	gorm.Model
	UserID int        `gorm:"not null;type:integer;uniqueIndex:idx_user_date,priority:1" json:"user_id"`
	Date   DateString `gorm:"not null;type:date;uniqueIndex:idx_user_date,priority:2" json:"date"`
	Mental int        `gorm:"not null;type:integer" json:"mental"`
	Diary  string     `gorm:"not null;type:text" json:"diary"`
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
		Date:   string(d.Date),
		Mental: mental,
		Diary:  d.Diary,
	}
}

// FromDomain converts the domain model to the persistence model.
func FromDomain(d *diary.Diary) *DiaryModel {
	return &DiaryModel{
		UserID: d.UserID,
		Date:   DateString(d.Date),
		Mental: d.Mental.GetValue(),
		Diary:  d.Diary,
	}
}
