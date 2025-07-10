package repositories

import (
	"errors"
	"strings"
	"tofunote-backend/domain/diary"
	"tofunote-backend/infra/db"

	"gorm.io/gorm"
)

type IDiaryRepository interface {
	FindAll() (*[]diary.Diary, error)
	FindByUserID(userID string) (*[]diary.Diary, error)
	FindByUserIDAndDate(userID string, date string) (*diary.Diary, error)
	FindByUserIDAndDateRange(userID string, startDate, endDate string) (*[]diary.Diary, error)
	Create(diary *diary.Diary) error
	Update(userID string, date string, diary *diary.Diary) error
	Delete(userID string, date string) error
	DeleteByUserID(userID string) error // 追加
	// FindByID(diaryId int) (*models.Diary, error)
}

type DiaryRepository struct {
	db *gorm.DB
}

func NewDiaryRepository(db *gorm.DB) IDiaryRepository {
	return &DiaryRepository{db: db}
}

func (r *DiaryRepository) FindAll() (*[]diary.Diary, error) {
	var diaryModels []db.DiaryModel
	if err := r.db.Find(&diaryModels).Error; err != nil {
		return nil, err
	}

	diaries := make([]diary.Diary, 0)
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return &diaries, nil
}

func (r *DiaryRepository) FindByUserID(userID string) (*[]diary.Diary, error) {
	var diaryModels []db.DiaryModel
	if err := r.db.Where("user_id = ?", userID).Find(&diaryModels).Error; err != nil {
		return nil, err
	}

	diaries := make([]diary.Diary, 0)
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return &diaries, nil
}

func (r *DiaryRepository) FindByUserIDAndDate(userID string, date string) (*diary.Diary, error) {
	var diaryModel db.DiaryModel
	if err := r.db.Where("user_id = ? AND date = ?", userID, date).First(&diaryModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("指定された日付の日記が見つかりません")
		}
		return nil, err
	}
	return diaryModel.ToDomain(), nil
}

func (r *DiaryRepository) FindByUserIDAndDateRange(userID string, startDate, endDate string) (*[]diary.Diary, error) {
	var diaryModels []db.DiaryModel
	if err := r.db.Where("user_id = ? AND date BETWEEN ? AND ?", userID, startDate, endDate).Find(&diaryModels).Error; err != nil {
		return nil, err
	}

	diaries := make([]diary.Diary, 0)
	for _, model := range diaryModels {
		diaries = append(diaries, *model.ToDomain())
	}
	return &diaries, nil
}

func (r *DiaryRepository) Create(diary *diary.Diary) error {
	model := db.FromDomain(diary)
	if err := r.db.Create(model).Error; err != nil {
		// 複合ユニークキー制約違反のエラーハンドリング
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("この日付の日記は既に作成されています")
		}
		return err
	}
	return nil
}

func (r *DiaryRepository) Update(userID string, date string, diary *diary.Diary) error {
	model := db.FromDomain(diary)
	result := r.db.Where("user_id = ? AND date = ?", userID, date).Updates(model)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("指定された日付の日記が見つかりません")
	}
	return nil
}

func (r *DiaryRepository) Delete(userID string, date string) error {
	result := r.db.Unscoped().Where("user_id = ? AND date = ?", userID, date).Delete(&db.DiaryModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("指定された日付の日記が見つかりません")
	}
	return nil
}

// 指定ユーザーの全日記を削除
func (r *DiaryRepository) DeleteByUserID(userID string) error {
	result := r.db.Unscoped().Where("user_id = ?", userID).Delete(&db.DiaryModel{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
