package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetPostDetailByID(ctx context.Context, id int) (*model.Post, error)
	GetPostsByUserID(ctx context.Context, userId int) ([]model.Post, error)
	UpdatePost(ctx context.Context, id int, post *model.Post) error
	DeletePost(ctx context.Context, id int) error
}

type PostRepositoryImpl struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &PostRepositoryImpl{db: db}
}

func (r *PostRepositoryImpl) CreatePost(ctx context.Context, post *model.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepositoryImpl) GetPostDetailByID(ctx context.Context, id int) (*model.Post, error) {
	var post model.Post
	if err := r.db.WithContext(ctx).First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepositoryImpl) DeletePost(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Post{}, id).Error
}

func (r *PostRepositoryImpl) UpdatePost(ctx context.Context, id int, post *model.Post) error {
	return r.db.WithContext(ctx).Model(&model.Post{}).Where("id = ?", id).Updates(post).Error
}

func (r *PostRepositoryImpl) GetPostsByUserID(ctx context.Context, userId int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.WithContext(ctx).Where("user_id", userId).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
