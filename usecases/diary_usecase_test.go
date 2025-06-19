package usecases

import (
	"emotra-backend/domain/diary"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// モックリポジトリ（IDiaryRepository の簡易実装）
type mockDiaryRepository struct {
	diaries *[]diary.Diary
	err     error
}

func (m *mockDiaryRepository) FindAll() (*[]diary.Diary, error) {
	return m.diaries, m.err
}

func (m *mockDiaryRepository) FindByUserIDAndDate(userID int, date string) (*diary.Diary, error) {
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

func (m *mockDiaryRepository) Create(diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryRepository) Update(userID int, date string, diary *diary.Diary) error {
	return m.err
}

func (m *mockDiaryRepository) Delete(userID int, date string) error {
	return m.err
}

// testDiaries はテスト用のダイアリーデータを定義します
var testDiaries = func() []diary.Diary {
	m5, _ := diary.NewMental(5)
	m3, _ := diary.NewMental(3)
	return []diary.Diary{
		{ID: 1, UserID: 101, Date: "2025-05-01", Mental: m5, Diary: "今日は良い一日だった"},
		{ID: 2, UserID: 102, Date: "2025-05-02", Mental: m3, Diary: "少し疲れた"},
	}
}()

func TestDiaryUsecase_FindAll(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func() *mockDiaryRepository
		expectedDiaries []diary.Diary
		expectedError   error
	}{
		{
			name: "正常系：リポジトリから取得した値をそのまま返す",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			expectedDiaries: testDiaries,
			expectedError:   nil,
		},
		{
			name: "正常系：空のリストを処理できる",
			setupMock: func() *mockDiaryRepository {
				emptyDiaries := make([]diary.Diary, 0)
				return &mockDiaryRepository{
					diaries: &emptyDiaries,
					err:     nil,
				}
			},
			expectedDiaries: []diary.Diary{},
			expectedError:   nil,
		},
		{
			name: "異常系：リポジトリがエラーを返した場合にそのまま返す",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: nil,
					err:     errors.New("DB接続失敗"),
				}
			},
			expectedDiaries: nil,
			expectedError:   errors.New("DB接続失敗"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			usecase := NewDiaryUsecase(mockRepo)

			result, err := usecase.FindAll()

			// エラーチェック
			if tt.expectedError != nil {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
					return
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラーが発生しました: %v", err)
			}

			if result == nil {
				t.Error("結果がnilです")
				return
			}

			// 結果の比較
			if diff := cmp.Diff(tt.expectedDiaries, *result); diff != "" {
				t.Errorf("期待値と実際の値が異なります:\n%s", diff)
			}
		})
	}
}

func TestDiaryUsecase_Create(t *testing.T) {
	m5, _ := diary.NewMental(5)

	tests := []struct {
		name          string
		setupMock     func() *mockDiaryRepository
		createDiary   diary.Diary
		expectedError error
	}{
		{
			name: "正常系：日記を作成できる",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: nil,
				}
			},
			createDiary: diary.Diary{
				UserID: 101,
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "新しい日記内容",
			},
			expectedError: nil,
		},
		{
			name: "異常系：複合ユニークキー制約違反",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("この日付の日記は既に作成されています"),
				}
			},
			createDiary: diary.Diary{
				UserID: 101,
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "重複する日記",
			},
			expectedError: errors.New("この日付の日記は既に作成されています"),
		},
		{
			name: "異常系：DBエラーが発生した場合",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("DB接続エラー"),
				}
			},
			createDiary: diary.Diary{
				UserID: 101,
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "DBエラー時の日記",
			},
			expectedError: errors.New("DB接続エラー"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			usecase := NewDiaryUsecase(mockRepo)

			err := usecase.Create(&tt.createDiary)

			// エラーチェック
			if tt.expectedError != nil {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
					return
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラーが発生しました: %v", err)
			}
		})
	}
}

func TestDiaryUsecase_Update(t *testing.T) {
	m5, _ := diary.NewMental(5)
	m7, _ := diary.NewMental(7)

	tests := []struct {
		name          string
		setupMock     func() *mockDiaryRepository
		userID        int
		date          string
		updateDiary   diary.Diary
		expectedError error
	}{
		{
			name: "正常系：日記を更新できる",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: nil,
				}
			},
			userID: 101,
			date:   "2025-05-01",
			updateDiary: diary.Diary{
				UserID: 101,
				Date:   "2025-05-01",
				Mental: m7,
				Diary:  "更新された日記内容",
			},
			expectedError: nil,
		},
		{
			name: "異常系：リポジトリがエラーを返した場合にそのまま返す",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("指定された日付の日記が見つかりません"),
				}
			},
			userID: 101,
			date:   "2025-05-01",
			updateDiary: diary.Diary{
				UserID: 101,
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "存在しない日記",
			},
			expectedError: errors.New("指定された日付の日記が見つかりません"),
		},
		{
			name: "異常系：DBエラーが発生した場合",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("DB接続エラー"),
				}
			},
			userID: 101,
			date:   "2025-05-01",
			updateDiary: diary.Diary{
				UserID: 101,
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "DBエラー時の日記",
			},
			expectedError: errors.New("DB接続エラー"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			usecase := NewDiaryUsecase(mockRepo)

			err := usecase.Update(tt.userID, tt.date, &tt.updateDiary)

			// エラーチェック
			if tt.expectedError != nil {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
					return
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラーが発生しました: %v", err)
			}
		})
	}
}

func TestDiaryUsecase_Delete(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func() *mockDiaryRepository
		userID        int
		date          string
		expectedError error
	}{
		{
			name: "正常系：日記を削除できる",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: nil,
				}
			},
			userID:        101,
			date:          "2025-05-01",
			expectedError: nil,
		},
		{
			name: "異常系：リポジトリがエラーを返した場合にそのまま返す",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("指定された日付の日記が見つかりません"),
				}
			},
			userID:        101,
			date:          "2025-05-01",
			expectedError: errors.New("指定された日付の日記が見つかりません"),
		},
		{
			name: "異常系：DBエラーが発生した場合",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					err: errors.New("DB接続エラー"),
				}
			},
			userID:        101,
			date:          "2025-05-01",
			expectedError: errors.New("DB接続エラー"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			usecase := NewDiaryUsecase(mockRepo)

			err := usecase.Delete(tt.userID, tt.date)

			// エラーチェック
			if tt.expectedError != nil {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
					return
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラーが発生しました: %v", err)
			}
		})
	}
}

func TestDiaryUsecase_FindByUserIDAndDate(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func() *mockDiaryRepository
		userID        int
		date          string
		expectedDiary *diary.Diary
		expectedError error
	}{
		{
			name: "正常系：指定されたuser_idとdateの日記を取得できる",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			userID:        101,
			date:          "2025-05-01",
			expectedDiary: &testDiaries[0],
			expectedError: nil,
		},
		{
			name: "異常系：指定されたuser_idとdateの日記が見つからない",
			setupMock: func() *mockDiaryRepository {
				diaries := make([]diary.Diary, len(testDiaries))
				copy(diaries, testDiaries)
				return &mockDiaryRepository{
					diaries: &diaries,
					err:     nil,
				}
			},
			userID:        999,
			date:          "2025-05-01",
			expectedDiary: nil,
			expectedError: errors.New("指定された日付の日記が見つかりません"),
		},
		{
			name: "異常系：リポジトリがエラーを返した場合にそのまま返す",
			setupMock: func() *mockDiaryRepository {
				return &mockDiaryRepository{
					diaries: nil,
					err:     errors.New("DB接続失敗"),
				}
			},
			userID:        101,
			date:          "2025-05-01",
			expectedDiary: nil,
			expectedError: errors.New("DB接続失敗"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := tt.setupMock()
			usecase := NewDiaryUsecase(mockRepo)

			result, err := usecase.FindByUserIDAndDate(tt.userID, tt.date)

			// エラーチェック
			if tt.expectedError != nil {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
					return
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("予期しないエラーが発生しました: %v", err)
			}

			if result == nil {
				t.Error("結果がnilです")
				return
			}

			// 結果の比較
			if diff := cmp.Diff(tt.expectedDiary, result); diff != "" {
				t.Errorf("期待値と実際の値が異なります:\n%s", diff)
			}
		})
	}
}
