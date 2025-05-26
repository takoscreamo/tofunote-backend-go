package repositories

import "emotra-backend/api/models"

type IDiaryRepository interface {
	FindAll() (*[]models.Diary, error)
}

type DiaryMemoryRepository struct {
	diaries []models.Diary
}

func NewDiaryMemoryRepository(diaries []models.Diary) IDiaryRepository {
	return &DiaryMemoryRepository{diaries: diaries}
}

func (r *DiaryMemoryRepository) FindAll() (*[]models.Diary, error) {
	return &r.diaries, nil
}
