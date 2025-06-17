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

func (m *mockDiaryRepository) Create(diary *diary.Diary) error {
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
