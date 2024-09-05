package model

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	gorm.Model
	ID           int          `gorm:"primary_key;column:id"`
	Name         string       `gorm:"column:name"`
	Universities []University `gorm:"foreignKey:LocationID"`
	CreatedAt    time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time    `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (l *Location) TableName() string {
	return "locations"
}
