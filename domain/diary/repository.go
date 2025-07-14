// DiaryRepositoryインターフェース: Diaryエンティティの永続化を抽象化する

package diary

import "context"

type DiaryRepository interface {
	FindAll(ctx context.Context) ([]Diary, error)
	FindByUserID(ctx context.Context, userID string) ([]Diary, error)
	FindByUserIDAndDate(ctx context.Context, userID string, date string) (*Diary, error)
	FindByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate string) ([]Diary, error)
	Create(ctx context.Context, diary *Diary) error
	Update(ctx context.Context, userID string, date string, diary *Diary) error
	Delete(ctx context.Context, userID string, date string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
