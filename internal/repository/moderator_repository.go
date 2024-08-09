package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type ModeratorRepository interface {
	CreateModerator(ctx context.Context, moderator *model.Moderator) error
	GetModeratorsByCommunityID(ctx context.Context, communityId int) ([]model.Moderator, error)
	DeleteModerator(ctx context.Context, id int) error
}

type ModeratorRepositoryImpl struct {
	db *gorm.DB
}

func NewModeratorRepository(db *gorm.DB) ModeratorRepository {
	return &ModeratorRepositoryImpl{
		db: db,
	}
}

func (r *ModeratorRepositoryImpl) CreateModerator(ctx context.Context, moderator *model.Moderator) error {
	return r.db.WithContext(ctx).Create(moderator).Error
}

func (r *ModeratorRepositoryImpl) GetModeratorsByCommunityID(ctx context.Context, communityId int) ([]model.Moderator, error) {
	var moderators []model.Moderator
	if err := r.db.WithContext(ctx).Where("community_id = ?", communityId).Error; err != nil {
		return nil, err
	}
	return moderators, nil
}

func (r *ModeratorRepositoryImpl) DeleteModerator(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Moderator{}, id).Error
}
