package usecases

import (
	"errors"
	"feelog-backend/domain/diary"
	"feelog-backend/domain/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDiaryRepo struct {
	deleteByUserIDErr error
}

func (m *mockDiaryRepo) DeleteByUserID(userID string) error {
	return m.deleteByUserIDErr
}
func (m *mockDiaryRepo) FindAll() (*[]diary.Diary, error)                   { return nil, nil }
func (m *mockDiaryRepo) FindByUserID(userID string) (*[]diary.Diary, error) { return nil, nil }
func (m *mockDiaryRepo) FindByUserIDAndDate(userID, date string) (*diary.Diary, error) {
	return nil, nil
}
func (m *mockDiaryRepo) FindByUserIDAndDateRange(userID, startDate, endDate string) (*[]diary.Diary, error) {
	return nil, nil
}
func (m *mockDiaryRepo) Create(d *diary.Diary) error                          { return nil }
func (m *mockDiaryRepo) Update(userID, date string, diary *diary.Diary) error { return nil }
func (m *mockDiaryRepo) Delete(userID, date string) error                     { return nil }

// 他のIDiaryRepositoryメソッドは未使用なので省略

type mockUserRepo struct {
	deleteByIDErr error
}

func (m *mockUserRepo) DeleteByID(id string) error {
	return m.deleteByIDErr
}
func (m *mockUserRepo) FindByProviderId(provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByRefreshToken(refreshToken string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) Create(u *user.User) error                                  { return nil }
func (m *mockUserRepo) FindByID(id string) (*user.User, error)                     { return nil, nil }
func (m *mockUserRepo) Update(u *user.User) error                                  { return nil }

// 他のuser.Repositoryメソッドは未使用なので省略

func TestUserWithdrawUsecase_Withdraw(t *testing.T) {
	tests := []struct {
		name              string
		diaryErr          error
		userErr           error
		expectError       bool
		expectErrorString string
	}{
		{
			name:        "正常系: 全削除成功",
			diaryErr:    nil,
			userErr:     nil,
			expectError: false,
		},
		{
			name:              "異常系: 日記削除失敗",
			diaryErr:          errors.New("diary error"),
			userErr:           nil,
			expectError:       true,
			expectErrorString: "diary error",
		},
		{
			name:              "異常系: ユーザー削除失敗",
			diaryErr:          nil,
			userErr:           errors.New("user error"),
			expectError:       true,
			expectErrorString: "user error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diaryRepo := &mockDiaryRepo{deleteByUserIDErr: tt.diaryErr}
			userRepo := &mockUserRepo{deleteByIDErr: tt.userErr}
			usecase := NewUserWithdrawUsecase(userRepo, diaryRepo)
			err := usecase.Withdraw("test-user")
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectErrorString, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
