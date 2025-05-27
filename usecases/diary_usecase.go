package usecases

import (
	"emotra-backend/api/models"
	"emotra-backend/repositories"
)

type IDiaryUsecase interface {
	FindAll() (*[]models.Diary, error)
}

type DiaryUsecase struct {
	repository repositories.IDiaryRepository
}

func NewDiaryUsecase(repository repositories.IDiaryRepository) IDiaryUsecase {
	return &DiaryUsecase{repository: repository}
}

func (s *DiaryUsecase) FindAll() (*[]models.Diary, error) {
	return s.repository.FindAll()
}
