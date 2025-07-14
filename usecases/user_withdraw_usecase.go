package usecases

import (
	"context"
	"tofunote-backend/domain/diary"
	"tofunote-backend/domain/user"
)

type UserWithdrawUsecase struct {
	UserRepository  user.Repository
	DiaryRepository diary.DiaryRepository
}

func NewUserWithdrawUsecase(userRepo user.Repository, diaryRepo diary.DiaryRepository) *UserWithdrawUsecase {
	return &UserWithdrawUsecase{
		UserRepository:  userRepo,
		DiaryRepository: diaryRepo,
	}
}

// Withdraw: 指定ユーザーの全日記とアカウントを削除
func (u *UserWithdrawUsecase) Withdraw(ctx context.Context, userID string) error {
	// 1. 日記全削除
	if err := u.DiaryRepository.DeleteByUserID(ctx, userID); err != nil {
		return err
	}
	// 2. ユーザー削除
	if err := u.UserRepository.DeleteByID(ctx, userID); err != nil {
		return err
	}
	return nil
}
