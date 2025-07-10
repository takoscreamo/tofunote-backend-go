package usecases

import (
	"tofunote-backend/domain/user"
	"tofunote-backend/repositories"
)

type UserWithdrawUsecase struct {
	UserRepository  user.Repository
	DiaryRepository repositories.IDiaryRepository
}

func NewUserWithdrawUsecase(userRepo user.Repository, diaryRepo repositories.IDiaryRepository) *UserWithdrawUsecase {
	return &UserWithdrawUsecase{
		UserRepository:  userRepo,
		DiaryRepository: diaryRepo,
	}
}

// Withdraw: 指定ユーザーの全日記とアカウントを削除
func (u *UserWithdrawUsecase) Withdraw(userID string) error {
	// 1. 日記全削除
	if err := u.DiaryRepository.DeleteByUserID(userID); err != nil {
		return err
	}
	// 2. ユーザー削除
	if err := u.UserRepository.DeleteByID(userID); err != nil {
		return err
	}
	return nil
}
