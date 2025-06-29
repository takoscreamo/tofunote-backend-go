package usecases

import (
	"feelog-backend/domain/diary"
	"feelog-backend/repositories"
)

type IDiaryUsecase interface {
	FindAll() (*[]diary.Diary, error)
	FindByUserID(userID int) (*[]diary.Diary, error)
	FindByUserIDAndDate(userID int, date string) (*diary.Diary, error)
	FindByUserIDAndDateRange(userID int, startDate, endDate string) (*[]diary.Diary, error)
	Create(diary *diary.Diary) error
	Update(userID int, date string, diary *diary.Diary) error
	Delete(userID int, date string) error
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

func (s *DiaryUsecase) FindByUserID(userID int) (*[]diary.Diary, error) {
	return s.repository.FindByUserID(userID)
}

func (s *DiaryUsecase) FindByUserIDAndDate(userID int, date string) (*diary.Diary, error) {
	return s.repository.FindByUserIDAndDate(userID, date)
}

func (s *DiaryUsecase) FindByUserIDAndDateRange(userID int, startDate, endDate string) (*[]diary.Diary, error) {
	return s.repository.FindByUserIDAndDateRange(userID, startDate, endDate)
}

func (s *DiaryUsecase) Create(diary *diary.Diary) error {
	return s.repository.Create(diary)
}

func (s *DiaryUsecase) Update(userID int, date string, diary *diary.Diary) error {
	return s.repository.Update(userID, date, diary)
}

func (s *DiaryUsecase) Delete(userID int, date string) error {
	return s.repository.Delete(userID, date)
}
