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

	var diaries []diary.Diary
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return &diaries, nil
}

// func (r *DiaryRepository) Create(ctx context.Context, d *diary.Diary) error {
// 	model := db.FromDomain(d)
// 	return r.db.WithContext(ctx).Create(model).Error
// }
