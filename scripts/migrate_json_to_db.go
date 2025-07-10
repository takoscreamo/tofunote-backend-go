package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"tofunote-backend/domain/diary"
	"tofunote-backend/infra"
	"tofunote-backend/repositories"
)

type DiaryEntry struct {
	Date    string `json:"date"`
	Mental  int    `json:"mental"`
	Content string `json:"content"`
}

type DiaryData struct {
	Entries []DiaryEntry `json:"entries"`
	Total   int          `json:"total"`
}

func main() {
	// データベース接続
	db := infra.SetupDB()
	diaryRepo := repositories.NewDiaryRepository(db)

	// JSONファイルを読み込み
	diaryData, err := loadJSONFile("infra/diary_datas/2025.json")
	if err != nil {
		log.Fatalf("JSONファイルの読み込みに失敗しました: %v", err)
	}

	log.Printf("JSONファイルから読み込まれた日記エントリ数: %d", len(diaryData.Entries))

	// データベースに移行
	userID := 1 // デフォルトのユーザーID
	successCount := 0
	errorCount := 0
	updatedCount := 0

	for _, entry := range diaryData.Entries {
		// メンタルスコアの正規化
		normalizedMental := normalizeMentalScore(entry.Mental)

		// メンタルスコアが0の場合はデフォルト値（5）を設定
		if normalizedMental == 0 {
			normalizedMental = 5
			log.Printf("メンタルスコアが見つからないため、デフォルト値5を設定 (日付: %s)", entry.Date)
		}

		// メンタルスコアを検証
		mental, err := diary.NewMental(normalizedMental)
		if err != nil {
			log.Printf("メンタルスコアが無効です (日付: %s, スコア: %d): %v", entry.Date, normalizedMental, err)
			errorCount++
			continue
		}

		// 日記エンティティを作成
		diaryEntry := &diary.Diary{
			UserID: userID,
			Date:   entry.Date,
			Mental: mental,
			Diary:  entry.Content,
		}

		// 既存の日記を確認
		existingDiary, err := diaryRepo.FindByUserIDAndDate(userID, entry.Date)
		if err != nil {
			// エラーが発生した場合（データが見つからない場合など）は新規作成を試行
			err = diaryRepo.Create(diaryEntry)
			if err != nil {
				log.Printf("日記の保存に失敗しました (日付: %s): %v", entry.Date, err)
				errorCount++
				continue
			}
			successCount++
			log.Printf("日記を新規作成しました: %s (メンタル: %d)", entry.Date, normalizedMental)
		} else {
			// 既存データがある場合は更新
			existingDiary.Mental = mental
			existingDiary.Diary = entry.Content

			err = diaryRepo.Update(userID, entry.Date, existingDiary)
			if err != nil {
				log.Printf("日記の更新に失敗しました (日付: %s): %v", entry.Date, err)
				errorCount++
				continue
			}
			updatedCount++
			log.Printf("日記を更新しました: %s (メンタル: %d)", entry.Date, normalizedMental)
		}
	}

	log.Printf("移行完了: 新規作成 %d件, 更新 %d件, 失敗 %d件", successCount, updatedCount, errorCount)
}

func loadJSONFile(filename string) (*DiaryData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("JSONファイルを開けませんでした: %v", err)
	}
	defer file.Close()

	var diaryData DiaryData
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&diaryData); err != nil {
		return nil, fmt.Errorf("JSONデコードに失敗しました: %v", err)
	}

	return &diaryData, nil
}

// メンタルスコアを正規化する（異常値を適切な範囲に収める）
func normalizeMentalScore(score int) int {
	if score <= 0 {
		return 0 // デフォルト値として扱う
	}
	if score > 10 {
		// 異常に高い値は10に制限
		return 10
	}
	return score
}
