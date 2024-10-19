package model

import (
	"time"

	"gorm.io/gorm"
)

type University struct {
	gorm.Model
	ID           int       `gorm:"primary_key;university_id"`
	Name         string    `gorm:"column:name"`
	Summary      string    `gorm:"column:summary"`
	LocationID   int       `gorm:"column:location_id"`
	Website      string    `gorm:"column:website"`
	Address      string    `gorm:"column:address"`
	TotalReviews *int      `gorm:"column:total_reviews"`
	Rating       *int      `gorm:"column:rating"`
	Reviews      []Review  `gorm:"foreignKey:UniversityID"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *University) TableName() string {
	return "universities"
}
