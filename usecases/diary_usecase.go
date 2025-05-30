package usecases

import (
	"emotra-backend/domain/diary"
	"emotra-backend/repositories"
)

type IDiaryUsecase interface {
	FindAll() (*[]diary.Diary, error)
}

type DiaryUsecase struct {
	repository repositories.IDiaryRepository
}

func NewDiaryUsecase(repository repositories.IDiaryRepository) IDiaryUsecase {
	return &DiaryUsecase{repository: repository}
}

func (s *DiaryUsecase) FindAll() (*[]diary.Diary, error) {
	return s.repository.FindAll()
}
