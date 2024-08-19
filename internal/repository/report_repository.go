package repository

import (
	"context"

	"github.com/temuka-api-service/internal/model"
	"gorm.io/gorm"
)

type ReportRepository interface {
	CreateReport(ctx context.Context, report *model.Report) error
	DeleteReport(ctx context.Context, id int) error
}

type ReportRepositoryImpl struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &ReportRepositoryImpl{
		db: db,
	}
}

func (r *ReportRepositoryImpl) CreateReport(ctx context.Context, report *model.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

func (r *ReportRepositoryImpl) DeleteReport(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Report{}, id).Error
}
