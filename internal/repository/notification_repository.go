package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *model.Notification) error
	GetNotificationsByUserID(ctx context.Context, userId int) ([]model.Notification, error)
}

type NotificationRepositoryImpl struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &NotificationRepositoryImpl{db: db}
}

func (r *NotificationRepositoryImpl) CreateNotification(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

func (r *NotificationRepositoryImpl) GetNotificationsByUserID(ctx context.Context, userId int) ([]model.Notification, error) {
	var notifications []model.Notification
	if err := r.db.WithContext(ctx).Where("user_id", userId).Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}
