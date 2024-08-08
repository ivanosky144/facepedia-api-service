package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type CommunityRepository interface {
	CreateCommunity(context context.Context, community *model.Community) error
	UpdateCommunity(context context.Context, id int, community *model.Community) error
	GetCommunityDetailByID(context context.Context, id int) (*model.Community, error)
	CheckMembership(ctx context.Context, communityID, userID int) (*model.CommunityMember, error)
	AddCommunityMember(ctx context.Context, member *model.CommunityMember) error
}

type CommunityRepositoryImpl struct {
	db *gorm.DB
}

func NewCommunityRepository(db *gorm.DB) CommunityRepository {
	return &CommunityRepositoryImpl{
		db: db,
	}
}

func (r *CommunityRepositoryImpl) CreateCommunity(ctx context.Context, community *model.Community) error {
	return r.db.WithContext(ctx).Create(community).Error
}

func (r *CommunityRepositoryImpl) UpdateCommunity(ctx context.Context, id int, community *model.Community) error {
	return r.db.WithContext(ctx).Model(&model.Community{}).Where("id = ?", id).Updates(community).Error
}

func (r *CommunityRepositoryImpl) GetCommunityDetailByID(ctx context.Context, id int) (*model.Community, error) {
	var community model.Community
	if err := r.db.WithContext(ctx).First(&community, id).Error; err != nil {
		return nil, err
	}
	return &community, nil
}

func (r *CommunityRepositoryImpl) AddCommunityMember(ctx context.Context, member *model.CommunityMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *CommunityRepositoryImpl) CheckMembership(ctx context.Context, communityID, userID int) (*model.CommunityMember, error) {
	var member model.CommunityMember
	if err := r.db.Where("community_id = ? AND user_id = ?", communityID, userID).First(&member).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}
