package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	// 日記ファイルを読み込み
	entries, err := parseDiaryFile("infra/diary_datas/2025.txt")
	if err != nil {
		log.Fatalf("日記ファイルの解析に失敗しました: %v", err)
	}

	// JSONデータ構造を作成
	diaryData := DiaryData{
		Entries: entries,
		Total:   len(entries),
	}

	// JSONファイルに出力
	outputFile := "infra/diary_datas/2025.json"
	err = writeJSONFile(outputFile, diaryData)
	if err != nil {
		log.Fatalf("JSONファイルの書き込みに失敗しました: %v", err)
	}

	log.Printf("JSONファイルを作成しました: %s", outputFile)
	log.Printf("変換された日記エントリ数: %d", len(entries))

	// 統計情報を表示
	mentalStats := make(map[int]int)
	for _, entry := range entries {
		mentalStats[entry.Mental]++
	}

	fmt.Println("\nメンタルスコアの統計:")
	fmt.Println("====================")
	for score := 1; score <= 10; score++ {
		count := mentalStats[score]
		fmt.Printf("スコア %d: %d件\n", score, count)
	}

	// メンタルスコアが0のエントリを表示
	fmt.Println("\nメンタルスコアが設定されていないエントリ:")
	fmt.Println("=====================================")
	for _, entry := range entries {
		if entry.Mental == 0 {
			fmt.Printf("日付: %s\n", entry.Date)
		}
	}
}

func writeJSONFile(filename string, data DiaryData) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ファイルを作成できませんでした: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("JSONエンコードに失敗しました: %v", err)
	}

	return nil
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
