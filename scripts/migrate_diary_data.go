package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"tofunote-backend/domain/diary"
	"tofunote-backend/infra"
	"tofunote-backend/repositories"
)

type DiaryEntry struct {
	Date    string
	Mental  int
	Content string
}

func main() {
	// データベース接続
	db := infra.SetupDB()
	diaryRepo := repositories.NewDiaryRepository(db)

	// 日記ファイルを読み込み
	entries, err := parseDiaryFile("infra/diary_datas/2025.txt")
	if err != nil {
		log.Fatalf("日記ファイルの解析に失敗しました: %v", err)
	}

	log.Printf("解析された日記エントリ数: %d", len(entries))

	// データベースに移行
	userID := 1 // デフォルトのユーザーID
	successCount := 0
	errorCount := 0

	for _, entry := range entries {
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

		// データベースに保存
		err = diaryRepo.Create(diaryEntry)
		if err != nil {
			log.Printf("日記の保存に失敗しました (日付: %s): %v", entry.Date, err)
			errorCount++
			continue
		}

		successCount++
		log.Printf("日記を保存しました: %s (メンタル: %d)", entry.Date, normalizedMental)
	}

	log.Printf("移行完了: 成功 %d件, 失敗 %d件", successCount, errorCount)
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

func parseDiaryFile(filename string) ([]DiaryEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ファイルを開けませんでした: %v", err)
	}
	defer file.Close()

	var entries []DiaryEntry
	scanner := bufio.NewScanner(file)

	// 日付とメンタルスコアを抽出する正規表現
	datePattern := regexp.MustCompile(`^(\d{1,2}/\d{1,2}[月火水木金土日]?)`)
	mentalPattern := regexp.MustCompile(`・メンタル(\d+)`)

	var currentEntry *DiaryEntry
	var currentContent strings.Builder

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 空行をスキップ
		if line == "" {
			continue
		}

		// 日付行を検出
		if dateMatch := datePattern.FindStringSubmatch(line); dateMatch != nil {
			// 前のエントリを保存
			if currentEntry != nil {
				currentEntry.Content = strings.TrimSpace(currentContent.String())
				entries = append(entries, *currentEntry)
			}

			// 新しいエントリを開始
			dateStr := dateMatch[1]
			formattedDate := formatDate(dateStr)

			currentEntry = &DiaryEntry{
				Date:    formattedDate,
				Mental:  0, // 後で設定
				Content: "",
			}
			currentContent.Reset()
		}

		// メンタルスコアを検出
		if mentalMatch := mentalPattern.FindStringSubmatch(line); mentalMatch != nil && currentEntry != nil {
			if mentalScore, err := strconv.Atoi(mentalMatch[1]); err == nil {
				currentEntry.Mental = mentalScore
			}
		}

		// コンテンツを追加
		if currentEntry != nil {
			if currentContent.Len() > 0 {
				currentContent.WriteString("\n")
			}
			currentContent.WriteString(line)
		}
	}

	// 最後のエントリを保存
	if currentEntry != nil {
		currentEntry.Content = strings.TrimSpace(currentContent.String())
		entries = append(entries, *currentEntry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ファイル読み込みエラー: %v", err)
	}

	return entries, nil
}

func formatDate(dateStr string) string {
	// "6/22日" のような形式を "2025-06-22" に変換
	re := regexp.MustCompile(`^(\d{1,2})/(\d{1,2})`)
	matches := re.FindStringSubmatch(dateStr)
	if len(matches) != 3 {
		return dateStr // 変換できない場合はそのまま返す
	}

	month, _ := strconv.Atoi(matches[1])
	day, _ := strconv.Atoi(matches[2])

	// 2025年として処理
	return fmt.Sprintf("2025-%02d-%02d", month, day)
}
