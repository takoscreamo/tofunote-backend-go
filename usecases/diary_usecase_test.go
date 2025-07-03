package usecases

import (
	"errors"
	"feelog-backend/domain/diary"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

// モックリポジトリ（IDiaryRepository の簡易実装）
type mockDiaryRepository struct {
	diaries *[]diary.Diary
	err     error
}

func (m *mockDiaryRepository) FindAll() (*[]diary.Diary, error) {
	return m.diaries, m.err
}

func (m *mockDiaryRepository) FindByUserID(userID string) (*[]diary.Diary, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.diaries == nil {
		emptySlice := make([]diary.Diary, 0)
		return &emptySlice, nil
	}
	// 特定のユーザーの日記のみをフィルタリング
	var userDiaries []diary.Diary
	for _, d := range *m.diaries {
		if d.UserID == userID {
			userDiaries = append(userDiaries, d)
		}
	}
	// 存在しないユーザーの場合も空配列を返す（nilではなく）
	return &userDiaries, nil
}

func (m *mockDiaryRepository) FindByUserIDAndDate(userID string, date string) (*diary.Diary, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.diaries == nil {
		return nil, errors.New("指定された日付の日記が見つかりません")
	}
	for _, d := range *m.diaries {
		if d.UserID == userID && d.Date == date {
			return &d, nil
		}
	}
	return nil, errors.New("指定された日付の日記が見つかりません")
}

func (m *mockDiaryRepository) FindByUserIDAndDateRange(userID string, startDate, endDate string) (*[]diary.Diary, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.diaries == nil {
		emptySlice := make([]diary.Diary, 0)
		return &emptySlice, nil
	}
	// 特定のユーザーの指定期間の日記のみをフィルタリング
	var userDiaries []diary.Diary
	for _, d := range *m.diaries {
		if d.UserID == userID && d.Date >= startDate && d.Date <= endDate {
			userDiaries = append(userDiaries, d)
		}
	}
	// 空の配列を返す（nilではなく）
	return &userDiaries, nil
}

func (m *mockDiaryRepository) Create(diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryRepository) Update(userID string, date string, diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryRepository) Delete(userID string, date string) error {
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
		expected  *[]diary.Diary
		hasError  bool
	}{
		{
			name: "正常系：全件取得できる",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			expected: &testDiaries,
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

			result, err := usecase.FindAll()

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
		expected  *[]diary.Diary
		hasError  bool
	}{
		{
			name:   "正常系：特定ユーザーの日記を取得できる",
			userID: "1",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			expected: &[]diary.Diary{
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
					diaries: &diaries,
					err:     nil,
				}
			},
			expected: &[]diary.Diary{},
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

			result, err := usecase.FindByUserID(tt.userID)

			if tt.hasError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				// より直接的な比較
				if result == nil && tt.expected == nil {
					// 両方ともnilの場合はOK
				} else if result == nil || tt.expected == nil {
					// 片方だけnilの場合はNG
					t.Errorf("FindByUserID() mismatch: result=%v, expected=%v", result, tt.expected)
				} else {
					// 両方ともnilでない場合は内容を比較
					if len(*result) != len(*tt.expected) {
						t.Errorf("FindByUserID() length mismatch: got=%d, expected=%d", len(*result), len(*tt.expected))
					} else {
						for i, r := range *result {
							e := (*tt.expected)[i]
							if r != e {
								t.Errorf("FindByUserID() item[%d] mismatch: got=%+v, expected=%+v", i, r, e)
							}
						}
					}
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
			name:   "正常系：指定されたユーザーと日付の日記を取得できる",
			userID: "1",
			date:   "2025-05-01",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			expected: &testDiaries[0],
			hasError: false,
		},
		{
			name:   "異常系：指定された日記が見つからない",
			userID: "1",
			date:   "2025-05-99",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
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

			result, err := usecase.FindByUserIDAndDate(tt.userID, tt.date)

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
	tests := []struct {
		name      string
		diary     *diary.Diary
		setupMock func() *mockDiaryRepository
		hasError  bool
	}{
		{
			name: "正常系：日記を作成できる",
			diary: &diary.Diary{
				UserID: "1",
				Date:   "2025-05-03",
				Mental: testDiaries[0].Mental,
				Diary:  "新しい日記",
			},
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: nil,
				}
			},
			hasError: false,
		},
		{
			name: "異常系：リポジトリがエラーを返す",
			diary: &diary.Diary{
				UserID: "1",
				Date:   "2025-05-03",
				Mental: testDiaries[0].Mental,
				Diary:  "新しい日記",
			},
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("DBエラー"),
				}
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			err := usecase.Create(tt.diary)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDiaryUsecase_Update(t *testing.T) {
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
			diary: &diary.Diary{
				UserID: "1",
				Date:   "2025-05-01",
				Mental: testDiaries[0].Mental,
				Diary:  "更新された日記",
			},
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: nil,
				}
			},
			hasError: false,
		},
		{
			name:   "異常系：リポジトリがエラーを返す",
			userID: "1",
			date:   "2025-05-01",
			diary: &diary.Diary{
				UserID: "1",
				Date:   "2025-05-01",
				Mental: testDiaries[0].Mental,
				Diary:  "更新された日記",
			},
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("DBエラー"),
				}
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			err := usecase.Update(tt.userID, tt.date, tt.diary)

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
					err: nil,
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
					err: errors.New("DBエラー"),
				}
			},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.setupMock()
			usecase := NewDiaryUsecase(mock)

			err := usecase.Delete(tt.userID, tt.date)

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
		expected  *[]diary.Diary
		hasError  bool
	}{
		{
			name:      "正常系：指定された期間の日記を取得できる",
			userID:    "1",
			startDate: "2025-05-01",
			endDate:   "2025-05-02",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			expected: &[]diary.Diary{
				{ID: "1", UserID: "1", Date: "2025-05-01", Mental: testDiaries[0].Mental, Diary: "今日は良い一日だった"},
				{ID: "2", UserID: "1", Date: "2025-05-02", Mental: testDiaries[1].Mental, Diary: "少し疲れた"},
			},
			hasError: false,
		},
		{
			name:      "正常系：指定された期間に日記がない場合は空配列を返す",
			userID:    "1",
			startDate: "2025-06-01",
			endDate:   "2025-06-30",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			expected: new([]diary.Diary),
			hasError: false,
		},
		{
			name:      "正常系：存在しないユーザーの場合は空配列を返す",
			userID:    "999",
			startDate: "2025-05-01",
			endDate:   "2025-05-02",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			expected: new([]diary.Diary),
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

			result, err := usecase.FindByUserIDAndDateRange(tt.userID, tt.startDate, tt.endDate)

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
