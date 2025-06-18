package repositories

import (
	"emotra-backend/domain/diary"
	"emotra-backend/infra/db"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type IDiaryRepository interface {
	FindAll() (*[]diary.Diary, error)
	Create(diary *diary.Diary) error
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

func (r *DiaryRepository) Create(diary *diary.Diary) error {
	model := db.FromDomain(diary)
	if err := r.db.Create(model).Error; err != nil {
		// 複合ユニークキー制約違反のエラーハンドリング
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("この日付の日記は既に作成されています")
		}
		return err
	}
	return nil
}
