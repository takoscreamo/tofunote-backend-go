// Diaryエンティティ: 日記のビジネスロジックを表現するモデル

package diary

type Diary struct {
	ID     string
	UserID string
	Date   string
	Mental Mental
	Diary  string
}
