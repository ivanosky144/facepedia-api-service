package model

import (
	"time"

	"gorm.io/gorm"
)

type MajorReview struct {
	gorm.Model
	ID        int       `gorm:"primary_key;column:id"`
	UserID    int       `gorm:"column:user_id"`
	MajorID   int       `gorm:"column:major_id"`
	Text      string    `gorm:column:text`
	Stars     int       `gorm:column:stars`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (m *MajorReview) TableName() string {
	return "major_reviews"
}
