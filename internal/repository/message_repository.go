package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type MessageRepository interface {
	CreateMessage(ctx context.Context, message *model.Message) error
	DeleteMessage(ctx context.Context, id int) error
}

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageRepositoryImpl(db *gorm.DB) MessageRepository {
	return &MessageRepositoryImpl{
		db: db,
	}
}

func (r *MessageRepositoryImpl) CreateMessage(ctx context.Context, message *model.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *MessageRepositoryImpl) DeleteMessage(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Message{}, id).Error
}
