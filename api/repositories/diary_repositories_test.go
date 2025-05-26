package repositories

import (
	"emotra-backend/api/models"
	"testing"
)

func TestDiaryMemoryRepository(t *testing.T) {
	t.Run("全てのダイアリーを取得できる", func(t *testing.T) {
		// テストデータを準備
		diaries := []models.Diary{
			{ID: 1, UserID: 101, Date: "2025-05-01", Mental: 5, Diary: "今日は楽しい一日だった。"},
			{ID: 2, UserID: 102, Date: "2025-05-02", Mental: 3, Diary: "少し疲れたけど頑張った。"},
		}

		// リポジトリを初期化
		repo := NewDiaryMemoryRepository(diaries)

		// FindAll メソッドを呼び出し
		result, err := repo.FindAll()
		if err != nil {
			t.Fatalf("エラーが発生しました: %v", err)
		}

		// 結果を検証
		if len(*result) != len(diaries) {
			t.Errorf("期待するダイアリー数: %d, 実際のダイアリー数: %d", len(diaries), len(*result))
		}

		for i, diary := range *result {
			if diary.ID != diaries[i].ID || diary.UserID != diaries[i].UserID || diary.Date != diaries[i].Date ||
				diary.Mental != diaries[i].Mental || diary.Diary != diaries[i].Diary {
				t.Errorf("期待するダイアリー: %+v, 実際のダイアリー: %+v", diaries[i], diary)
			}
		}
	})

	t.Run("空のダイアリーリストを処理できる", func(t *testing.T) {
		// 空のリポジトリを初期化
		repo := NewDiaryMemoryRepository([]models.Diary{})

		// FindAll メソッドを呼び出し
		result, err := repo.FindAll()
		if err != nil {
			t.Fatalf("エラーが発生しました: %v", err)
		}

		// 結果を検証
		if len(*result) != 0 {
			t.Errorf("期待するダイアリー数: 0, 実際のダイアリー数: %d", len(*result))
		}
	})
}
