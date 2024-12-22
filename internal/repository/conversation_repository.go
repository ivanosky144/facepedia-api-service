package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	CreateConversation(ctx context.Context, conversation *model.Conversation) error
	GetConversationsByUserID(ctx context.Context, userID int) ([]model.Conversation, error)
	DeleteConversation(ctx context.Context, id int) error
	GetConversationDetailByID(ctx context.Context, id int) (*model.Conversation, error)
	AddParticipant(ctx context.Context, participant *model.Participant) error
	AddMessage(ctx context.Context, message *model.Message) error
}

type ConversationRepositoryImpl struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &ConversationRepositoryImpl{
		db: db,
	}
}

func (r *ConversationRepositoryImpl) CreateConversation(ctx context.Context, conversation *model.Conversation) error {
	return r.db.WithContext(ctx).Create(conversation).Error
}

func (r *ConversationRepositoryImpl) GetConversationsByUserID(ctx context.Context, userID int) ([]model.Conversation, error) {
	var conversations []model.Conversation
	if err := r.db.Where("user_id = ?", userID).Find(&conversations).Error; err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *ConversationRepositoryImpl) DeleteConversation(ctx context.Context, id int) error {
	return r.db.Delete(&model.Conversation{}, id).Error
}

func (r *ConversationRepositoryImpl) AddMessage(ctx context.Context, message *model.Message) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *ConversationRepositoryImpl) GetConversationDetailByID(ctx context.Context, id int) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := r.db.WithContext(ctx).First(&conversation, id).Error; err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepositoryImpl) AddParticipant(ctx context.Context, participant *model.Participant) error {
	return r.db.WithContext(ctx).Create(participant).Error

}
