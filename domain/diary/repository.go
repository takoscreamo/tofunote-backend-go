// DiaryRepositoryインターフェース: Diaryエンティティの永続化を抽象化する

package diary

import "context"

type DiaryRepository interface {
	FindAll(ctx context.Context) ([]Diary, error)
	Create(ctx context.Context, diary *Diary) error
}
