package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error)
	DeleteComment(ctx context.Context, commentID int) error
	GetRepliesByParentID(ctx context.Context, parentID int) ([]model.Comment, error)
}

type CommentRepositoryImpl struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &CommentRepositoryImpl{
		db: db,
	}
}

func (r *CommentRepositoryImpl) CreateComment(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *CommentRepositoryImpl) GetCommentsByPostID(ctx context.Context, postID int) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepositoryImpl) DeleteComment(ctx context.Context, commentID int) error {
	return r.db.Delete(&model.Comment{}, commentID).Error
}

func (r *CommentRepositoryImpl) GetRepliesByParentID(ctx context.Context, parentID int) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.Where("parent_id = ?", parentID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
