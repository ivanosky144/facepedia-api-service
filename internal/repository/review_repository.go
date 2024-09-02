package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *model.Review) error
	DeleteReview(ctx context.Context, id int) error
	GetReviewsByUniversityID(ctx context.Context, universityID int) ([]model.Review, error)
}

type ReviewRepositoryImpl struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &ReviewRepositoryImpl{
		db: db,
	}
}

func (r *ReviewRepositoryImpl) CreateReview(ctx context.Context, review *model.Review) error {
	return r.db.WithContext(ctx).Create(&review).Error
}

func (r *ReviewRepositoryImpl) DeleteReview(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Review{}, id).Error
}

func (r *ReviewRepositoryImpl) GetReviewsByUniversityID(ctx context.Context, universityID int) ([]model.Review, error) {
	var reviews []model.Review
	if err := r.db.WithContext(ctx).Where("university_id = ?", universityID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
