package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type CommunityRepository interface {
	CreateCommunity(context context.Context, community *model.Community) error
	CheckCommunityNameAvailability(ctx context.Context, name string) bool
	UpdateCommunity(context context.Context, id int, community *model.Community) error
	GetCommunities(context context.Context) ([]model.Community, error)
	GetUserJoinedCommunities(context context.Context, userID int) ([]model.Community, error)
	GetCommunityDetailByID(context context.Context, id int) (*model.Community, error)
	CheckMembership(ctx context.Context, communityID, userID int) (*model.CommunityMember, error)
	AddCommunityMember(ctx context.Context, member *model.CommunityMember) error
	GetCommunityPosts(ctx context.Context, id int, filters map[string]interface{}) ([]model.CommunityPost, error)
	UpdateCommunityPostsCount(context context.Context, id int) error
	UpdateCommunityMembersCount(context context.Context, id int) error
	DeleteCommunity(context context.Context, id int) error
	GetCommunityDetailBySlug(ctx context.Context, slug string) (*model.Community, error)
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

func (r *CommunityRepositoryImpl) CheckCommunityNameAvailability(ctx context.Context, name string) bool {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Community{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false
	}
	return count == 0
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

func (r *CommunityRepositoryImpl) GetCommunities(context context.Context) ([]model.Community, error) {
	var communities []model.Community
	if err := r.db.WithContext(context).Find(&communities).Error; err != nil {
		return nil, err
	}
	return communities, nil
}

func (r *CommunityRepositoryImpl) DeleteCommunity(context context.Context, id int) error {
	return r.db.WithContext(context).Delete(&model.Community{}, id).Error
}

func (r *CommunityRepositoryImpl) AddCommunityMember(ctx context.Context, member *model.CommunityMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *CommunityRepositoryImpl) UpdateCommunityPostsCount(context context.Context, id int) error {
	return r.db.WithContext(context).Model(&model.Community{}).Where("id = ?", id).
		Update("posts_count", gorm.Expr("posts_count + 1")).Error
}

func (r *CommunityRepositoryImpl) UpdateCommunityMembersCount(context context.Context, id int) error {
	return r.db.WithContext(context).Model(&model.Community{}).Where("id = ?", id).
		Update("members_count", gorm.Expr("members_count + 1")).Error
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

func (r *CommunityRepositoryImpl) GetCommunityPosts(ctx context.Context, communityID int, filters map[string]interface{}) ([]model.CommunityPost, error) {
	var communityPosts []model.CommunityPost

	data := r.db.WithContext(ctx).Where("community_id = ?", communityID)

	for key, val := range filters {
		if key == "sort" || key == "sort_by" {
			continue
		}
		data = data.Where(key+" = ?", val)
	}

	sortBy, sortByExists := filters["sort_by"].(string)
	sortOrder, sortOrderExists := filters["sort"].(string)

	if sortByExists && sortOrderExists {
		data = data.Order(sortBy + " " + sortOrder)
	} else if sortByExists {
		data = data.Order(sortBy + "asc")
	} else {
		data = data.Order("created_at desc")
	}

	if err := data.Find(&communityPosts).Error; err != nil {
		return nil, err
	}

	return communityPosts, nil
}

func (r *CommunityRepositoryImpl) GetUserJoinedCommunities(context context.Context, userID int) ([]model.Community, error) {
	var communities []model.Community

	query := `
		SELECT c.*
		FROM community_members cm
		INNER JOIN communities c ON cm.community_id = c.id
		WHERE cm.user_id = ? AND cm.banned = false
	`

	if err := r.db.WithContext(context).Raw(query, userID).Scan(&communities).Error; err != nil {
		return nil, err
	}

	return communities, nil
}

func (r *CommunityRepositoryImpl) GetCommunityDetailBySlug(ctx context.Context, slug string) (*model.Community, error) {
	var community model.Community
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&community).Error; err != nil {
		return nil, err
	}
	return &community, nil
}
