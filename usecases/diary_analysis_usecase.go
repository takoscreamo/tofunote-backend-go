package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"emotra-backend/repositories"
)

type DiaryAnalysisUsecase struct {
	DiaryRepository repositories.IDiaryRepository
}

func NewDiaryAnalysisUsecase(diaryRepository repositories.IDiaryRepository) *DiaryAnalysisUsecase {
	return &DiaryAnalysisUsecase{
		DiaryRepository: diaryRepository,
	}
}

type AnalysisRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type AnalysisResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// AnalyzeDiary は日記の内容を分析する
func (u *DiaryAnalysisUsecase) AnalyzeAllDiaries() (string, error) {
	diaries, err := u.DiaryRepository.FindAll()
	if err != nil {
		return "", err
	}

	// 日記の内容を結合
	var diaryContents []string
	for _, diary := range *diaries {
		// 各フィールドを結合
		entry := strings.Join([]string{
			"ID: " + strconv.Itoa(diary.ID),
			"UserID: " + strconv.Itoa(diary.UserID),
			"Date: " + diary.Date,
			"Mental: " + strconv.Itoa(diary.Mental.GetValue()),
			"Diary: " + diary.Diary,
		}, "\n")
		diaryContents = append(diaryContents, entry)
	}
	combinedContent := strings.Join(diaryContents, "\n\n")

	// APIリクエストデータの作成
	requestBody := AnalysisRequest{
		Model: "deepseek/deepseek-r1-0528-qwen3-8b:free",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "system", Content: "あなたはユーザーの日記とメンタルスコア（1〜10）をもとに、感情の傾向を分析し、やさしく前向きなアドバイスを行うメンタルサポートAIです。\n\nユーザーのメンタルスコアは1〜10の10段階で記録されており、1が最も調子が悪く、10が最も調子が良いことを表します。\n\nスコアと日記の内容を組み合わせて、感情の傾向を読み取り、簡潔に100文字以内で説明してください。"},
			{Role: "user", Content: "以下はユーザーの日記とメンタルスコアです。\n\n" + combinedContent + "\n\nこの内容を分析して、感情の傾向を読み取り、わかりやすく丁寧に説明してください。"},
		},
	}

	// JSONエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	// Authorization トークンを環境変数から取得
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", errors.New("APIキーが設定されていません")
	}

	// APIリクエストの送信
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("APIリクエストが失敗しました: " + resp.Status)
	}

	// レスポンスのデコード
	var analysisResponse AnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&analysisResponse); err != nil {
		return "", err
	}

	// 分析結果の取得
	if len(analysisResponse.Choices) == 0 {
		return "", errors.New("分析結果が空です")
	}

	return analysisResponse.Choices[0].Message.Content, nil
}
