package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type UniversityRepository interface {
	CreateUniversity(ctx context.Context, university *model.University) error
	UpdateUniversity(ctx context.Context, id int, university *model.University) error
	GetUniversities(ctx context.Context) ([]model.University, error)
	DeleteUniversity(ctx context.Context, id int) error
	GetUniversityDetailByID(ctx context.Context, id int) (*model.University, error)
	GetUniversityDetailBySlug(ctx context.Context, slug string) (*model.University, error)
}

type UniversityRepositoryImpl struct {
	db *gorm.DB
}

func NewUniversityRepository(db *gorm.DB) UniversityRepository {
	return &UniversityRepositoryImpl{
		db: db,
	}
}

func (r *UniversityRepositoryImpl) CreateUniversity(ctx context.Context, university *model.University) error {
	return r.db.WithContext(ctx).Create(university).Error
}

func (r *UniversityRepositoryImpl) DeleteUniversity(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.University{}, id).Error
}

func (r *UniversityRepositoryImpl) UpdateUniversity(ctx context.Context, id int, post *model.University) error {
	return r.db.WithContext(ctx).Model(&model.University{}).Where("id = ?", id).Updates(post).Error
}

func (r *UniversityRepositoryImpl) GetUniversities(ctx context.Context) ([]model.University, error) {
	var universities []model.University
	if err := r.db.WithContext(ctx).Find(&universities).Error; err != nil {
		return nil, err
	}
	return universities, nil
}

func (r *UniversityRepositoryImpl) GetUniversityDetailByID(ctx context.Context, id int) (*model.University, error) {
	var university model.University
	if err := r.db.WithContext(ctx).First(&university, id).Error; err != nil {
		return nil, err
	}
	return &university, nil
}

func (r *UniversityRepositoryImpl) GetUniversityDetailBySlug(ctx context.Context, slug string) (*model.University, error) {
	var university model.University
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&university).Error; err != nil {
		return nil, err
	}
	return &university, nil
}
