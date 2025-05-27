package repositories

import (
	"emotra-backend/api/models"

	"gorm.io/gorm"
)

type IDiaryRepository interface {
	FindAll() (*[]models.Diary, error)
	// FindByID(diaryId int) (*models.Diary, error)
}

type DiaryRepository struct {
	db *gorm.DB
}

func NewDiaryRepository(db *gorm.DB) IDiaryRepository {
	return &DiaryRepository{db: db}
}

func (r *DiaryRepository) FindAll() (*[]models.Diary, error) {
	var diaries []models.Diary
	result := r.db.Find(&diaries)
	if result.Error != nil {
		return nil, result.Error
	}
	return &diaries, nil
}
