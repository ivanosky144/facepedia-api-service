package models

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	ID        int       `gorm:"primary_key;column:id"`
	Title     string    `gorm:"column:title"`
	CreatorID string    `gorm:"column:creator_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Conversation) TableName() string {
	return "conversations"
}
