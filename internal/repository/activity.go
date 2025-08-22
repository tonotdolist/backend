package repository

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
	"tonotdolist/common"
	"tonotdolist/internal/model"
)

type ActivityRepository interface {
	GetNUserActivity(ctx context.Context, userId string, n int, offset int) ([]*model.Activity, error)
	GetUserActivityInRange(ctx context.Context, userId string, start time.Time, end time.Time) ([]*model.Activity, error)
	CreateActivity(ctx context.Context, activity *model.Activity) error
	UpdateActivity(ctx context.Context, activity *model.Activity) error
	DeleteActivity(ctx context.Context, activityId string, userId string) error
}

type activityRepository struct {
	*Repository
}

func NewActivityRepository(repository *Repository) ActivityRepository {
	return &activityRepository{
		Repository: repository,
	}
}

func (r *activityRepository) GetNUserActivity(ctx context.Context, userId string, n int, offset int) ([]*model.Activity, error) {
	var activities []*model.Activity
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Order("created_at DESC").Offset(offset).Limit(n).Find(activities).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, fmt.Errorf("error fetching n user activities from db: %v", err)
	}

	return activities, nil
}

func (r *activityRepository) GetUserActivityInRange(ctx context.Context, userId string, start time.Time, end time.Time) ([]*model.Activity, error) {
	var activities []*model.Activity
	err := r.db.WithContext(ctx).Where("user_id = ? AND date BETWEEN ? AND ?", userId, start, end).Find(&activities).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, fmt.Errorf("error fetching user activities from range from db: %v", err)
	}

	return activities, nil
}

func (r *activityRepository) CreateActivity(ctx context.Context, activity *model.Activity) error {
	err := r.db.WithContext(ctx).Create(&activity).Error
	if err != nil {
		return fmt.Errorf("error creating activity: %v", err)
	}

	return nil
}

func (r *activityRepository) UpdateActivity(ctx context.Context, activity *model.Activity) error {
	res := r.db.WithContext(ctx).Where("activity_id = ? AND user_id = ?", activity.ActivityId, activity.UserId).Updates(&activity)

	if err := res.Error; err != nil {
		return fmt.Errorf("error updating activity: %v", err)
	}

	if res.RowsAffected == 0 {
		return common.ErrNotFound
	}

	return nil
}

func (r *activityRepository) DeleteActivity(ctx context.Context, activityId string, userId string) error {
	res := r.db.WithContext(ctx).Where("activity_id = ? AND user_id = ?", activityId, userId).Delete(&model.Activity{})
	if err := res.Error; err != nil {
		return fmt.Errorf("error deleting activity: %v", err)
	}

	if res.RowsAffected == 0 {
		return common.ErrNotFound
	}

	return nil
}
