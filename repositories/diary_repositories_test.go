package repositories

import (
	"context"
	"errors"
	"testing"
	"time"
	"tofunote-backend/domain/diary"

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
		{ID: "1", UserID: "101", Date: "2025-05-01", Mental: m5, Diary: "今日は楽しい一日だった。"},
		{ID: "2", UserID: "102", Date: "2025-05-02", Mental: m3, Diary: "少し疲れたけど頑張った。"},
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
						AddRow("1", "101", "2025-05-01", 5, "今日は楽しい一日だった。", time.Now(), time.Now(), nil).
						AddRow("2", "102", "2025-05-02", 3, "少し疲れたけど頑張った。", time.Now(), time.Now(), nil),
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
			result, err := repo.FindAll(context.Background())

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

			if diff := cmp.Diff(tt.expectedDiaries, result); diff != "" {
				t.Errorf("期待値と実際の値が異なります:\n%s", diff)
			}

			verifyMockExpectations(t, mock)
		})
	}
}

func TestCreate(t *testing.T) {
	m5, _ := diary.NewMental(5)

	tests := []struct {
		name         string
		setupMock    func(sqlmock.Sqlmock)
		createDiary  diary.Diary
		expectError  bool
		errorMessage string
	}{
		{
			name: "正常系：日記を作成できる",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "diaries"`).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
				mock.ExpectCommit()
			},
			createDiary: diary.Diary{
				UserID: "101",
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "新しい日記内容",
			},
			expectError: false,
		},
		{
			name: "異常系：複合ユニークキー制約違反",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "diaries"`).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
				mock.ExpectRollback()
			},
			createDiary: diary.Diary{
				UserID: "101",
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "重複する日記",
			},
			expectError:  true,
			errorMessage: "この日付の日記は既に作成されています",
		},
		{
			name: "異常系：DBエラー",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "diaries"`).
					WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			createDiary: diary.Diary{
				UserID: "101",
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "DBエラー時の日記",
			},
			expectError:  true,
			errorMessage: "DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupTestDB(t)
			tt.setupMock(mock)
			repo := NewDiaryRepository(gormDB)

			err := repo.Create(context.Background(), &tt.createDiary)
			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
			}
			verifyMockExpectations(t, mock)
		})
	}
}

func TestUpdate(t *testing.T) {
	m5, _ := diary.NewMental(5)
	tests := []struct {
		name         string
		setupMock    func(sqlmock.Sqlmock)
		userID       string
		date         string
		updateDiary  diary.Diary
		expectError  bool
		errorMessage string
	}{
		{
			name: "正常系：日記を更新できる",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "diaries"`).
					WithArgs(sqlmock.AnyArg(), "101", "2025-05-01", 5, "更新内容", "101", "2025-05-01").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			userID: "101",
			date:   "2025-05-01",
			updateDiary: diary.Diary{
				UserID: "101",
				Date:   "2025-05-01",
				Mental: m5,
				Diary:  "更新内容",
			},
			expectError: false,
		},
		{
			name: "異常系：該当データなし",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "diaries"`).
					WithArgs(sqlmock.AnyArg(), "101", "2025-05-01", 5, "更新内容", "101", "2025-05-01").
					WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			userID:       "101",
			date:         "2025-05-01",
			updateDiary:  diary.Diary{UserID: "101", Date: "2025-05-01", Mental: m5, Diary: "更新内容"},
			expectError:  true,
			errorMessage: "指定された日付の日記が見つかりません",
		},
		{
			name: "異常系：DBエラー",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "diaries"`).
					WithArgs(sqlmock.AnyArg(), "101", "2025-05-01", 5, "更新内容", "101", "2025-05-01").
					WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			userID:       "101",
			date:         "2025-05-01",
			updateDiary:  diary.Diary{UserID: "101", Date: "2025-05-01", Mental: m5, Diary: "更新内容"},
			expectError:  true,
			errorMessage: "DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupTestDB(t)
			tt.setupMock(mock)
			repo := NewDiaryRepository(gormDB)

			err := repo.Update(context.Background(), tt.userID, tt.date, &tt.updateDiary)
			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
			}
			verifyMockExpectations(t, mock)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name         string
		setupMock    func(sqlmock.Sqlmock)
		userID       string
		date         string
		expectError  bool
		errorMessage string
	}{
		{
			name: "正常系：日記を削除できる",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "diaries"`).
					WithArgs("101", "2025-05-01").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			userID:      "101",
			date:        "2025-05-01",
			expectError: false,
		},
		{
			name: "異常系：該当データなし",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "diaries"`).
					WithArgs("101", "2025-05-01").
					WillReturnResult(sqlmock.NewResult(1, 0))
				mock.ExpectCommit()
			},
			userID:       "101",
			date:         "2025-05-01",
			expectError:  true,
			errorMessage: "指定された日付の日記が見つかりません",
		},
		{
			name: "異常系：DBエラー",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "diaries"`).
					WithArgs("101", "2025-05-01").
					WillReturnError(errors.New("DB error"))
				mock.ExpectRollback()
			},
			userID:       "101",
			date:         "2025-05-01",
			expectError:  true,
			errorMessage: "DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupTestDB(t)
			tt.setupMock(mock)
			repo := NewDiaryRepository(gormDB)

			err := repo.Delete(context.Background(), tt.userID, tt.date)
			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
			}
			verifyMockExpectations(t, mock)
		})
	}
}

func TestFindByUserIDAndDate(t *testing.T) {
	m5, _ := diary.NewMental(5)
	expectedDiary := diary.Diary{
		ID:     "1",
		UserID: "101",
		Date:   "2025-05-01",
		Mental: m5,
		Diary:  "テスト日記",
	}

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		userID        string
		date          string
		expectedDiary *diary.Diary
		expectError   bool
		errorMessage  string
	}{
		{
			name: "正常系：指定されたuser_idとdateの日記を取得できる",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "date", "mental", "diary", "created_at", "updated_at", "deleted_at"}).
					AddRow("1", "101", "2025-05-01", 5, "テスト日記", time.Now(), time.Now(), nil)
				mock.ExpectQuery(`SELECT \* FROM "diaries"`).
					WithArgs("101", "2025-05-01", 1).
					WillReturnRows(rows)
			},
			userID:        "101",
			date:          "2025-05-01",
			expectedDiary: &expectedDiary,
			expectError:   false,
		},
		{
			name: "異常系：指定されたuser_idとdateの日記が見つからない",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "diaries"`).
					WithArgs("999", "2025-05-01", 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			userID:        "999",
			date:          "2025-05-01",
			expectedDiary: nil,
			expectError:   true,
			errorMessage:  "指定された日付の日記が見つかりません",
		},
		{
			name: "異常系：DBエラー",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "diaries"`).
					WithArgs("101", "2025-05-01", 1).
					WillReturnError(errors.New("DB error"))
			},
			userID:        "101",
			date:          "2025-05-01",
			expectedDiary: nil,
			expectError:   true,
			errorMessage:  "DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupTestDB(t)
			tt.setupMock(mock)
			repo := NewDiaryRepository(gormDB)

			result, err := repo.FindByUserIDAndDate(context.Background(), tt.userID, tt.date)
			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if result == nil {
					t.Error("結果がnilです")
				} else {
					if result.UserID != tt.expectedDiary.UserID {
						t.Errorf("期待するUserID: %v, 実際のUserID: %v", tt.expectedDiary.UserID, result.UserID)
					}
					if result.Date != tt.expectedDiary.Date {
						t.Errorf("期待するDate: %v, 実際のDate: %v", tt.expectedDiary.Date, result.Date)
					}
					if result.Mental != tt.expectedDiary.Mental {
						t.Errorf("期待するMental: %v, 実際のMental: %v", tt.expectedDiary.Mental, result.Mental)
					}
					if result.Diary != tt.expectedDiary.Diary {
						t.Errorf("期待するDiary: %v, 実際のDiary: %v", tt.expectedDiary.Diary, result.Diary)
					}
				}
			}
			verifyMockExpectations(t, mock)
		})
	}
}

func TestFindByUserIDAndDateRange(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(sqlmock.Sqlmock)
		userID          string
		startDate       string
		endDate         string
		expectedDiaries []diary.Diary
		expectError     bool
		errorMessage    string
	}{
		{
			name: "正常系：指定された期間の日記を取得できる",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "date", "mental", "diary", "created_at", "updated_at", "deleted_at"}).
					AddRow("1", "101", "2025-05-01", 5, "今日は楽しい一日だった。", time.Now(), time.Now(), nil).
					AddRow("2", "101", "2025-05-02", 3, "少し疲れたけど頑張った。", time.Now(), time.Now(), nil)
				mock.ExpectQuery(`SELECT \* FROM "diaries"`).
					WithArgs("101", "2025-05-01", "2025-05-31").
					WillReturnRows(rows)
			},
			userID:    "101",
			startDate: "2025-05-01",
			endDate:   "2025-05-31",
			expectedDiaries: []diary.Diary{
				{ID: "1", UserID: "101", Date: "2025-05-01", Mental: testDiaries[0].Mental, Diary: "今日は楽しい一日だった。"},
				{ID: "2", UserID: "101", Date: "2025-05-02", Mental: testDiaries[1].Mental, Diary: "少し疲れたけど頑張った。"},
			},
			expectError: false,
		},
		{
			name: "正常系：指定された期間に日記がない場合は空配列を返す",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "date", "mental", "diary", "created_at", "updated_at", "deleted_at"})
				mock.ExpectQuery(`SELECT \* FROM "diaries"`).
					WithArgs("101", "2025-06-01", "2025-06-30").
					WillReturnRows(rows)
			},
			userID:          "101",
			startDate:       "2025-06-01",
			endDate:         "2025-06-30",
			expectedDiaries: []diary.Diary{},
			expectError:     false,
		},
		{
			name: "異常系：DBエラー",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "diaries"`).
					WithArgs("101", "2025-05-01", "2025-05-31").
					WillReturnError(errors.New("DB error"))
			},
			userID:          "101",
			startDate:       "2025-05-01",
			endDate:         "2025-05-31",
			expectedDiaries: nil,
			expectError:     true,
			errorMessage:    "DB error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gormDB, mock := setupTestDB(t)
			tt.setupMock(mock)
			repo := NewDiaryRepository(gormDB)

			result, err := repo.FindByUserIDAndDateRange(context.Background(), tt.userID, tt.startDate, tt.endDate)

			if tt.expectError {
				if err == nil {
					t.Error("エラーが期待されていましたが、発生しませんでした")
				} else if tt.errorMessage != "" && err.Error() != tt.errorMessage {
					t.Errorf("期待するエラー: %v, 実際のエラー: %v", tt.errorMessage, err)
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

			if diff := cmp.Diff(tt.expectedDiaries, result); diff != "" {
				t.Errorf("期待値と実際の値が異なります:\n%s", diff)
			}

			verifyMockExpectations(t, mock)
		})
	}
}
