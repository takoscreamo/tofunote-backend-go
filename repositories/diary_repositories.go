package repositories

import (
	"emotra-backend/domain/diary"
	"emotra-backend/infra/db"

	"gorm.io/gorm"
)

type IDiaryRepository interface {
	FindAll() (*[]diary.Diary, error)
	// FindByID(diaryId int) (*models.Diary, error)
}

type DiaryRepository struct {
	db *gorm.DB
}

func NewDiaryRepository(db *gorm.DB) IDiaryRepository {
	return &DiaryRepository{db: db}
}

func (r *DiaryRepository) FindAll() (*[]diary.Diary, error) {
	var diaryModels []db.DiaryModel
	if err := r.db.Find(&diaryModels).Error; err != nil {
		return nil, err
	}

	diaries := make([]diary.Diary, 0)
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return &diaries, nil
}
