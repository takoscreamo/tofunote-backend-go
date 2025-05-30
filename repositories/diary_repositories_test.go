package repositories

import (
	"emotra-backend/domain/diary"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDiaryRepository(t *testing.T) {
	t.Run("全てのダイアリーを取得できる", func(t *testing.T) {
		// モックDBの作成
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock の作成に失敗しました: %v", err)
		}
		defer db.Close()

		// GORM DBインスタンスの作成
		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		if err != nil {
			t.Fatalf("GORM DB の作成に失敗しました: %v", err)
		}

		// モックの期待値を設定
		mock.ExpectQuery(`SELECT \* FROM "diaries"`).WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "date", "mental", "diary"}).
				AddRow(1, 101, "2025-05-01", 5, "今日は楽しい一日だった。").
				AddRow(2, 102, "2025-05-02", 3, "少し疲れたけど頑張った。"),
		)

		// リポジトリを初期化
		repo := NewDiaryRepository(gormDB)

		// FindAll メソッドを呼び出し
		result, err := repo.FindAll()
		if err != nil {
			t.Fatalf("エラーが発生しました: %v", err)
		}

		// 結果を検証
		if len(*result) != 2 {
			t.Errorf("期待するダイアリー数: 2, 実際のダイアリー数: %d", len(*result))
		}

		expectedDiaries := []diary.Diary{
			{ID: 1, UserID: 101, Date: "2025-05-01", Mental: diary.NewMental(5), Diary: "今日は楽しい一日だった。"},
			{ID: 2, UserID: 102, Date: "2025-05-02", Mental: diary.NewMental(3), Diary: "少し疲れたけど頑張った。"},
		}

		for i, diary := range *result {
			if diary != expectedDiaries[i] {
				t.Errorf("期待するダイアリー: %+v, 実際のダイアリー: %+v", expectedDiaries[i], diary)
			}
		}

		// モックの検証
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("モックの期待値が満たされていません: %v", err)
		}
	})

	t.Run("空のダイアリーリストを処理できる", func(t *testing.T) {
		// モックDBの作成
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("sqlmock の作成に失敗しました: %v", err)
		}
		defer db.Close()

		// GORM DBインスタンスの作成
		gormDB, err := gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), &gorm.Config{})
		if err != nil {
			t.Fatalf("GORM DB の作成に失敗しました: %v", err)
		}

		// モックの期待値を設定
		mock.ExpectQuery(`SELECT \* FROM "diaries"`).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "date", "mental", "diary"}))

		// リポジトリを初期化
		repo := NewDiaryRepository(gormDB)

		// FindAll メソッドを呼び出し
		result, err := repo.FindAll()
		if err != nil {
			t.Fatalf("エラーが発生しました: %v", err)
		}

		// 結果を検証
		if len(*result) != 0 {
			t.Errorf("期待するダイアリー数: 0, 実際のダイアリー数: %d", len(*result))
		}

		// モックの検証
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("モックの期待値が満たされていません: %v", err)
		}
	})
}
