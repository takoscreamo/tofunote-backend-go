package main

import (
	"emotra-backend/infra"
	"emotra-backend/repositories"
	"fmt"
	"log"
)

func main() {
	// データベース接続
	db := infra.SetupDB()
	diaryRepo := repositories.NewDiaryRepository(db)

	// 全ユーザーの日記を取得
	diaries, err := diaryRepo.FindAll()
	if err != nil {
		log.Fatalf("日記の取得に失敗しました: %v", err)
	}

	fmt.Printf("データベース内の日記総数: %d件\n", len(*diaries))
	fmt.Println()

	// 最新の10件を表示
	fmt.Println("最新の10件の日記:")
	fmt.Println("==================")

	count := 0
	for i := len(*diaries) - 1; i >= 0 && count < 10; i-- {
		diary := (*diaries)[i]
		fmt.Printf("日付: %s, メンタル: %d\n", diary.Date, diary.Mental.Value())
		fmt.Printf("内容: %s\n", truncateString(diary.Diary, 100))
		fmt.Println("---")
		count++
	}

	// メンタルスコアの統計
	fmt.Println("メンタルスコアの統計:")
	fmt.Println("====================")

	mentalStats := make(map[int]int)
	for _, d := range *diaries {
		mentalStats[d.Mental.Value()]++
	}

	for score := 1; score <= 10; score++ {
		count := mentalStats[score]
		fmt.Printf("スコア %d: %d件\n", score, count)
	}
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
