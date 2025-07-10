package usecases

import (
	"tofunote-backend/domain/diary"
	"tofunote-backend/repositories"
)

type IDiaryUsecase interface {
	FindAll() (*[]diary.Diary, error)
	FindByUserID(userID string) (*[]diary.Diary, error)
	FindByUserIDAndDate(userID string, date string) (*diary.Diary, error)
	FindByUserIDAndDateRange(userID string, startDate, endDate string) (*[]diary.Diary, error)
	Create(diary *diary.Diary) error
	Update(userID string, date string, diary *diary.Diary) error
	Delete(userID string, date string) error
	DeleteByUserID(userID string) error
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

func (s *DiaryUsecase) FindByUserID(userID string) (*[]diary.Diary, error) {
	return s.repository.FindByUserID(userID)
}

func (s *DiaryUsecase) FindByUserIDAndDate(userID string, date string) (*diary.Diary, error) {
	return s.repository.FindByUserIDAndDate(userID, date)
}

func (s *DiaryUsecase) FindByUserIDAndDateRange(userID string, startDate, endDate string) (*[]diary.Diary, error) {
	return s.repository.FindByUserIDAndDateRange(userID, startDate, endDate)
}

func (s *DiaryUsecase) Create(diary *diary.Diary) error {
	return s.repository.Create(diary)
}

func (s *DiaryUsecase) Update(userID string, date string, diary *diary.Diary) error {
	return s.repository.Update(userID, date, diary)
}

func (s *DiaryUsecase) Delete(userID string, date string) error {
	return s.repository.Delete(userID, date)
}

func (s *DiaryUsecase) DeleteByUserID(userID string) error {
	return s.repository.DeleteByUserID(userID)
}
