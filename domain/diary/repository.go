// DiaryRepositoryインターフェース: Diaryエンティティの永続化を抽象化する

package diary

import "context"

type DiaryRepository interface {
	FindAll(ctx context.Context) ([]Diary, error)
	Create(ctx context.Context, diary *Diary) error
	Update(ctx context.Context, userID int, date string, diary *Diary) error
}
