package repositories

import (
	"emotra-backend/domain/diary"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// setupTestDB はテスト用のDBとモックをセットアップします
func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock の作成に失敗しました: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("GORM DB の作成に失敗しました: %v", err)
	}

	return gormDB, mock
}

// verifyMockExpectations はモックの期待値が満たされているかを検証します
func verifyMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("モックの期待値が満たされていません: %v", err)
	}
}

// testDiary はテスト用のダイアリーデータを定義します
var testDiaries = func() []diary.Diary {
	m5, _ := diary.NewMental(5)
	m3, _ := diary.NewMental(3)
	return []diary.Diary{
		{ID: 1, UserID: 101, Date: "2025-05-01", Mental: m5, Diary: "今日は楽しい一日だった。"},
		{ID: 2, UserID: 102, Date: "2025-05-02", Mental: m3, Diary: "少し疲れたけど頑張った。"},
	}
}()

func TestFindAll(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(sqlmock.Sqlmock)
		expectedDiaries []diary.Diary
		expectedError   bool
	}{
		{
			name: "全てのダイアリーを取得できる",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`(?i)SELECT \* FROM "diaries" WHERE "diaries"."deleted_at" IS NULL`).WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "date", "mental", "diary", "created_at", "updated_at", "deleted_at"}).
						AddRow(1, 101, "2025-05-01", 5, "今日は楽しい一日だった。", time.Now(), time.Now(), nil).
						AddRow(2, 102, "2025-05-02", 3, "少し疲れたけど頑張った。", time.Now(), time.Now(), nil),
				)
			},
			expectedDiaries: testDiaries,
			expectedError:   false,
		},
		{
			name: "空のダイアリーリストを処理できる",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`(?i)SELECT \* FROM "diaries" WHERE "diaries"."deleted_at" IS NULL`).WillReturnRows(
					sqlmock.NewRows([]string{"id", "user_id", "date", "mental", "diary", "created_at", "updated_at", "deleted_at"}),
				)
			},
			expectedDiaries: []diary.Diary{},
			expectedError:   false,
		},
		{
			name: "DBエラーが発生する",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`(?i)SELECT \* FROM "diaries" WHERE "diaries"."deleted_at" IS NULL`).WillReturnError(errors.New("DB error"))
			},
			expectedDiaries: nil,
			expectedError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupTestDB(t)

			tt.setupMock(mock)

			repo := NewDiaryRepository(gormDB)
			result, err := repo.FindAll()

			if tt.expectedError {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
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

			if diff := cmp.Diff(tt.expectedDiaries, *result); diff != "" {
				t.Errorf("期待値と実際の値が異なります:\n%s", diff)
			}

			verifyMockExpectations(t, mock)
		})
	}
}
