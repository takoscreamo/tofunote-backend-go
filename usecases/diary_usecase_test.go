package usecases

import (
	"emotra-backend/domain/diary"
	"errors"
	"testing"
)

// モックリポジトリ（IDiaryRepository の簡易実装）
type mockDiaryRepository struct {
	diaries *[]diary.Diary
	err     error
}

func (m *mockDiaryRepository) FindAll() (*[]diary.Diary, error) {
	return m.diaries, m.err
}

func TestDiaryUsecase_FindAll(t *testing.T) {
	t.Run("リポジトリから取得した値をそのまま返す", func(t *testing.T) {
		expected := &[]diary.Diary{
			{ID: 1, UserID: 101, Date: "2025-05-01", Mental: diary.NewMental(5), Diary: "今日は良い一日だった"},
			{ID: 2, UserID: 102, Date: "2025-05-02", Mental: diary.NewMental(3), Diary: "少し疲れた"},
		}

		mockRepo := &mockDiaryRepository{diaries: expected}
		usecase := NewDiaryUsecase(mockRepo)

		result, err := usecase.FindAll()
		if err != nil {
			t.Fatalf("エラーが返されました: %v", err)
		}

		if result != expected {
			t.Errorf("リポジトリの返却値と異なります。期待: %p, 実際: %p", expected, result)
		}
	})

	t.Run("リポジトリがエラーを返した場合にそのまま返す", func(t *testing.T) {
		expectedErr := errors.New("DB接続失敗")

		mockRepo := &mockDiaryRepository{
			diaries: nil,
			err:     expectedErr,
		}
		usecase := NewDiaryUsecase(mockRepo)

		_, err := usecase.FindAll()
		if err != expectedErr {
			t.Errorf("エラーの透過失敗。期待: %v, 実際: %v", expectedErr, err)
		}
	})
}
