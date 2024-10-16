package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type LocationRepository interface {
	AddLocation(ctx context.Context, location *model.Location) error
	UpdateLocation(ctx context.Context, id int, location *model.Location) error
	GetLocations(ctx context.Context) ([]model.Location, error)
	DeleteLocation(ctx context.Context, id int) error
	GetLocationById(ctx context.Context, id int) (*model.Location, error)
}

type LocationRepositoryImpl struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &LocationRepositoryImpl{
		db: db,
	}
}

func (r *LocationRepositoryImpl) AddLocation(ctx context.Context, location *model.Location) error {
	return r.db.WithContext(ctx).Create(location).Error
}

func (r *LocationRepositoryImpl) DeleteLocation(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Location{}, id).Error
}

func (r *LocationRepositoryImpl) UpdateLocation(ctx context.Context, id int, location *model.Location) error {
	return r.db.WithContext(ctx).Model(&model.Location{}).Where("id = ?", id).Updates(location).Error
}

func (r *LocationRepositoryImpl) GetLocations(ctx context.Context) ([]model.Location, error) {
	var locations []model.Location
	if err := r.db.WithContext(ctx).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *LocationRepositoryImpl) GetLocationById(ctx context.Context, id int) (*model.Location, error) {
	var location model.Location
	if err := r.db.WithContext(ctx).First(&location, id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}
