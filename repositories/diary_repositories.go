package repositories

import (
	"context"
	"errors"
	"strings"
	"tofunote-backend/domain/diary"
	"tofunote-backend/infra/db"

	"github.com/cmackenzie1/go-uuid"
	"gorm.io/gorm"
)

type DiaryRepository struct {
	db *gorm.DB
}

func NewDiaryRepository(db *gorm.DB) diary.DiaryRepository {
	return &DiaryRepository{db: db}
}

func (r *DiaryRepository) FindAll(ctx context.Context) ([]diary.Diary, error) {
	var diaryModels []db.DiaryModel
	if err := r.db.WithContext(ctx).Find(&diaryModels).Error; err != nil {
		return nil, err
	}

	diaries := make([]diary.Diary, 0, len(diaryModels))
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return diaries, nil
}

func (r *DiaryRepository) FindByUserID(ctx context.Context, userID string) ([]diary.Diary, error) {
	var diaryModels []db.DiaryModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&diaryModels).Error; err != nil {
		return nil, err
	}

	diaries := make([]diary.Diary, 0, len(diaryModels))
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return diaries, nil
}

func (r *DiaryRepository) FindByUserIDAndDate(ctx context.Context, userID string, date string) (*diary.Diary, error) {
	var diaryModel db.DiaryModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND date = ?", userID, date).First(&diaryModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("指定された日付の日記が見つかりません")
		}
		return nil, err
	}
	return diaryModel.ToDomain(), nil
}

func (r *DiaryRepository) FindByUserIDAndDateRange(ctx context.Context, userID string, startDate, endDate string) ([]diary.Diary, error) {
	var diaryModels []db.DiaryModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&diaryModels).Error; err != nil {
		return nil, err
	}

	diaries := make([]diary.Diary, 0, len(diaryModels))
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return diaries, nil
}

func (r *DiaryRepository) Create(ctx context.Context, diary *diary.Diary) error {
	if diary.ID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		diary.ID = id.String()
	}
	model := db.FromDomain(diary)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		// 複合ユニークキー制約違反のエラーハンドリング
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("この日付の日記は既に作成されています")
		}
		return err
	}
	return nil
}

func (r *DiaryRepository) Update(ctx context.Context, userID string, date string, diary *diary.Diary) error {
	model := db.FromDomain(diary)
	result := r.db.WithContext(ctx).Where("user_id = ? AND date = ?", userID, date).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("指定された日付の日記が見つかりません")
	}
	return nil
}

func (r *DiaryRepository) Delete(ctx context.Context, userID string, date string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("user_id = ? AND date = ?", userID, date).Delete(&db.DiaryModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("指定された日付の日記が見つかりません")
	}
	return nil
}

// 指定ユーザーの全日記を削除
func (r *DiaryRepository) DeleteByUserID(ctx context.Context, userID string) error {
	result := r.db.WithContext(ctx).Unscoped().Where("user_id = ?", userID).Delete(&db.DiaryModel{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
