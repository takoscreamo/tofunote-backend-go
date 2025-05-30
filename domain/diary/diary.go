// Diaryエンティティ: 日記のビジネスロジックを表現するモデル

package diary

type Diary struct {
	ID     int
	UserID int
	Date   string
	Mental Mental
	Diary  string
}
