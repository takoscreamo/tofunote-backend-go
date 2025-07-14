package usecases

import (
	"context"
	"tofunote-backend/domain/diary"
)

type IDiaryUsecase interface {
	FindAll(ctx context.Context) ([]diary.Diary, error)
	FindByUserID(ctx context.Context, userID string) ([]diary.Diary, error)
	FindByUserIDAndDate(ctx context.Context, userID string, date string) (*diary.Diary, error)
	FindByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate string) ([]diary.Diary, error)
	Create(ctx context.Context, diary *diary.Diary) error
	Update(ctx context.Context, userID string, date string, diary *diary.Diary) error
	Delete(ctx context.Context, userID string, date string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

type DiaryUsecase struct {
	repository diary.DiaryRepository
}

func NewDiaryUsecase(repository diary.DiaryRepository) IDiaryUsecase {
	return &DiaryUsecase{repository: repository}
}

func (s *DiaryUsecase) FindAll(ctx context.Context) ([]diary.Diary, error) {
	return s.repository.FindAll(ctx)
}

func (s *DiaryUsecase) FindByUserID(ctx context.Context, userID string) ([]diary.Diary, error) {
	return s.repository.FindByUserID(ctx, userID)
}

func (s *DiaryUsecase) FindByUserIDAndDate(ctx context.Context, userID string, date string) (*diary.Diary, error) {
	return s.repository.FindByUserIDAndDate(ctx, userID, date)
}

func (s *DiaryUsecase) FindByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate string) ([]diary.Diary, error) {
	return s.repository.FindByUserIDAndDateRange(ctx, userID, startDate, endDate)
}

func (s *DiaryUsecase) Create(ctx context.Context, diary *diary.Diary) error {
	return s.repository.Create(ctx, diary)
}

func (s *DiaryUsecase) Update(ctx context.Context, userID string, date string, diary *diary.Diary) error {
	return s.repository.Update(ctx, userID, date, diary)
}

func (s *DiaryUsecase) Delete(ctx context.Context, userID string, date string) error {
	return s.repository.Delete(ctx, userID, date)
}

func (s *DiaryUsecase) DeleteByUserID(ctx context.Context, userID string) error {
	return s.repository.DeleteByUserID(ctx, userID)
}
