package usecases

import (
	"emotra-backend/domain/diary"
	"emotra-backend/repositories"
)

type IDiaryUsecase interface {
	FindAll() (*[]diary.Diary, error)
	Create(diary *diary.Diary) error
	Update(userID int, date string, diary *diary.Diary) error
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

func (s *DiaryUsecase) Create(diary *diary.Diary) error {
	return s.repository.Create(diary)
}

func (s *DiaryUsecase) Update(userID int, date string, diary *diary.Diary) error {
	return s.repository.Update(userID, date, diary)
}
