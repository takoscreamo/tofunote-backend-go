package usecases

import (
	"context"
	"errors"
	"testing"
	"tofunote-backend/domain/diary"
	"tofunote-backend/domain/user"

	"github.com/stretchr/testify/assert"
)

type mockDiaryRepo struct {
	deleteByUserIDErr error
}

func (m *mockDiaryRepo) DeleteByUserID(ctx context.Context, userID string) error {
	return m.deleteByUserIDErr
}
func (m *mockDiaryRepo) FindAll(ctx context.Context) ([]diary.Diary, error) { return nil, nil }
func (m *mockDiaryRepo) FindByUserID(ctx context.Context, userID string) ([]diary.Diary, error) {
	return nil, nil
}
func (m *mockDiaryRepo) FindByUserIDAndDate(ctx context.Context, userID, date string) (*diary.Diary, error) {
	return nil, nil
}
func (m *mockDiaryRepo) FindByUserIDAndDateRange(ctx context.Context, userID, startDate, endDate string) ([]diary.Diary, error) {
	return nil, nil
}
func (m *mockDiaryRepo) Create(ctx context.Context, d *diary.Diary) error { return nil }
func (m *mockDiaryRepo) Update(ctx context.Context, userID, date string, diary *diary.Diary) error {
	return nil
}
func (m *mockDiaryRepo) Delete(ctx context.Context, userID, date string) error { return nil }

// 他のIDiaryRepositoryメソッドは未使用なので省略

type mockUserRepo struct {
	deleteByIDErr error
}

func (m *mockUserRepo) DeleteByID(ctx context.Context, id string) error {
	return m.deleteByIDErr
}
func (m *mockUserRepo) FindByProviderId(ctx context.Context, provider, providerId string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindByRefreshToken(ctx context.Context, refreshToken string) (*user.User, error) {
	return nil, nil
}
func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error              { return nil }
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*user.User, error) { return nil, nil }
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error              { return nil }

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
			err := usecase.Withdraw(context.Background(), "test-user")
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectErrorString, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
