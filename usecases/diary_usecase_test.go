package usecases

import (
	"context"
	"errors"
	"testing"
	"tofunote-backend/domain/diary"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// モックリポジトリ（domain/diary.DiaryRepository の簡易実装）
type mockDiaryRepository struct {
	diaries []diary.Diary
	err     error
}

func (m *mockDiaryRepository) FindAll(ctx context.Context) ([]diary.Diary, error) {
	return m.diaries, m.err
}

func (m *mockDiaryRepository) FindByUserID(ctx context.Context, userID string) ([]diary.Diary, error) {
	if m.err != nil {
		return nil, m.err
	}
	var userDiaries []diary.Diary
	if m.diaries != nil {
		for _, d := range m.diaries {
			if d.UserID == userID {
				userDiaries = append(userDiaries, d)
			}
		}
	}
	if userDiaries == nil {
		userDiaries = make([]diary.Diary, 0)
	}
	return userDiaries, nil
}

func (m *mockDiaryRepository) FindByUserIDAndDate(ctx context.Context, userID string, date string) (*diary.Diary, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, d := range m.diaries {
		if d.UserID == userID && d.Date == date {
			return &d, nil
		}
	}
	return nil, errors.New("指定された日付の日記が見つかりません")
}

func (m *mockDiaryRepository) FindByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate string) ([]diary.Diary, error) {
	if m.err != nil {
		return nil, m.err
	}
	var userDiaries []diary.Diary
	if m.diaries != nil {
		for _, d := range m.diaries {
			if d.UserID == userID && d.Date >= startDate && d.Date <= endDate {
				userDiaries = append(userDiaries, d)
			}
		}
	}
	if userDiaries == nil {
		userDiaries = make([]diary.Diary, 0)
	}
	return userDiaries, nil
}

func (m *mockDiaryRepository) Create(ctx context.Context, diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryRepository) Update(ctx context.Context, userID string, date string, diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryRepository) Delete(ctx context.Context, userID string, date string) error {
	return m.err
}

func (m *mockDiaryRepository) DeleteByUserID(ctx context.Context, userID string) error {
	return m.err
}

// testDiaries はテスト用のダイアリーデータを定義します
var testDiaries = func() []diary.Diary {
	m5, _ := diary.NewMental(5)
	m3, _ := diary.NewMental(3)
	return []diary.Diary{
		{ID: "1", UserID: "1", Date: "2025-05-01", Mental: m5, Diary: "今日は良い一日だった"},
		{ID: "2", UserID: "1", Date: "2025-05-02", Mental: m3, Diary: "少し疲れた"},
		{ID: "3", UserID: "2", Date: "2025-05-01", Mental: m5, Diary: "別のユーザーの日記"},
	}
}()

func TestDiaryUsecase_FindAll(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func() *mockDiaryRepository
		expected  []diary.Diary
		hasError  bool
	}{
		{
			name: "正常系：全件取得できる",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: diaries,
					err:     nil,
				}
			},
			expected: testDiaries,
			hasError: false,
		},
		{
			name: "異常系：リポジトリがエラーを返す",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: nil,
					err:     errors.New("DBエラー"),
				}
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			result, err := usecase.FindAll(context.Background())

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("FindAll() mismatch (-expected +got):\n%s", diff)
				}
			}
		})
	}
}

func TestDiaryUsecase_FindByUserID(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		setupMock func() *mockDiaryRepository
		expected  []diary.Diary
		hasError  bool
	}{
		{
			name:   "正常系：特定ユーザーの日記を取得できる",
			userID: "1",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: diaries,
					err:     nil,
				}
			},
			expected: []diary.Diary{
				{ID: "1", UserID: "1", Date: "2025-05-01", Mental: testDiaries[0].Mental, Diary: "今日は良い一日だった"},
				{ID: "2", UserID: "1", Date: "2025-05-02", Mental: testDiaries[1].Mental, Diary: "少し疲れた"},
			},
			hasError: false,
		},
		{
			name:   "正常系：存在しないユーザーの場合は空配列を返す",
			userID: "999",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: diaries,
					err:     nil,
				}
			},
			expected: []diary.Diary{},
			hasError: false,
		},
		{
			name:   "異常系：リポジトリがエラーを返す",
			userID: "1",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: nil,
					err:     errors.New("DBエラー"),
				}
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			result, err := usecase.FindByUserID(context.Background(), tt.userID)

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("FindByUserID() mismatch (-expected +got):\n%s", diff)
				}
			}
		})
	}
}

func TestDiaryUsecase_FindByUserIDAndDate(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		date      string
		setupMock func() *mockDiaryRepository
		expected  *diary.Diary
		hasError  bool
	}{
		{
			name:   "正常系：特定ユーザーの特定日付の日記を取得できる",
			userID: "1",
			date:   "2025-05-01",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: diaries,
					err:     nil,
				}
			},
			expected: &testDiaries[0],
			hasError: false,
		},
		{
			name:   "異常系：存在しない日記の場合はエラーを返す",
			userID: "1",
			date:   "2025-05-03",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: diaries,
					err:     nil,
				}
			},
			expected: nil,
			hasError: true,
		},
		{
			name:   "異常系：リポジトリがエラーを返す",
			userID: "1",
			date:   "2025-05-01",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: nil,
					err:     errors.New("DBエラー"),
				}
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			result, err := usecase.FindByUserIDAndDate(context.Background(), tt.userID, tt.date)

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("FindByUserIDAndDate() mismatch (-expected +got):\n%s", diff)
				}
			}
		})
	}
}

func TestDiaryUsecase_Create(t *testing.T) {
	m5, _ := diary.NewMental(5)
	newDiary := &diary.Diary{
		UserID: "1",
		Date:   "2025-05-03",
		Mental: m5,
		Diary:  "新しい日記",
	}

	tests := []struct {
		name      string
		diary     *diary.Diary
		setupMock func() *mockDiaryRepository
		hasError  bool
	}{
		{
			name:  "正常系：日記を作成できる",
			diary: newDiary,
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: testDiaries,
					err:     nil,
				}
			},
			hasError: false,
		},
		{
			name:  "異常系：リポジトリがエラーを返す",
			diary: newDiary,
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: testDiaries,
					err:     errors.New("DBエラー"),
				}
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			err := usecase.Create(context.Background(), tt.diary)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDiaryUsecase_Update(t *testing.T) {
	m4, _ := diary.NewMental(4)
	updateDiary := &diary.Diary{
		UserID: "1",
		Date:   "2025-05-01",
		Mental: m4,
		Diary:  "更新された日記",
	}

	tests := []struct {
		name      string
		userID    string
		date      string
		diary     *diary.Diary
		setupMock func() *mockDiaryRepository
		hasError  bool
	}{
		{
			name:   "正常系：日記を更新できる",
			userID: "1",
			date:   "2025-05-01",
			diary:  updateDiary,
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: testDiaries,
					err:     nil,
				}
			},
			hasError: false,
		},
		{
			name:   "異常系：リポジトリがエラーを返す",
			userID: "1",
			date:   "2025-05-01",
			diary:  updateDiary,
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: testDiaries,
					err:     errors.New("DBエラー"),
				}
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			err := usecase.Update(context.Background(), tt.userID, tt.date, tt.diary)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDiaryUsecase_Delete(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		date      string
		setupMock func() *mockDiaryRepository
		hasError  bool
	}{
		{
			name:   "正常系：日記を削除できる",
			userID: "1",
			date:   "2025-05-01",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: testDiaries,
					err:     nil,
				}
			},
			hasError: false,
		},
		{
			name:   "異常系：リポジトリがエラーを返す",
			userID: "1",
			date:   "2025-05-01",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: testDiaries,
					err:     errors.New("DBエラー"),
				}
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			err := usecase.Delete(context.Background(), tt.userID, tt.date)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDiaryUsecase_FindByUserIDAndDateRange(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		startDate string
		endDate   string
		setupMock func() *mockDiaryRepository
		expected  []diary.Diary
		hasError  bool
	}{
		{
			name:      "正常系：指定期間の日記を取得できる",
			userID:    "1",
			startDate: "2025-05-01",
			endDate:   "2025-05-02",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: diaries,
					err:     nil,
				}
			},
			expected: []diary.Diary{
				{ID: "1", UserID: "1", Date: "2025-05-01", Mental: testDiaries[0].Mental, Diary: "今日は良い一日だった"},
				{ID: "2", UserID: "1", Date: "2025-05-02", Mental: testDiaries[1].Mental, Diary: "少し疲れた"},
			},
			hasError: false,
		},
		{
			name:      "正常系：期間外の場合は空配列を返す",
			userID:    "1",
			startDate: "2025-05-10",
			endDate:   "2025-05-15",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: diaries,
					err:     nil,
				}
			},
			expected: []diary.Diary{},
			hasError: false,
		},
		{
			name:      "異常系：リポジトリがエラーを返す",
			userID:    "1",
			startDate: "2025-05-01",
			endDate:   "2025-05-02",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: nil,
					err:     errors.New("DBエラー"),
				}
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			result, err := usecase.FindByUserIDAndDateRange(context.Background(), tt.userID, tt.startDate, tt.endDate)

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if diff := cmp.Diff(tt.expected, result); diff != "" {
					t.Errorf("FindByUserIDAndDateRange() mismatch (-expected +got):\n%s", diff)
				}
			}
		})
	}
}
