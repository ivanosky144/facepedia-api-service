package models

import (
	"time"

	"gorm.io/gorm"
)

type CommunityPost struct {
	gorm.Model
	ID          int       `gorm:"primary_key;column:id"`
	PostID      string    `gorm:"foreign_key;column:post_id"`
	Mark        string    `gorm:"column:mark"`
	Topic       string    `gorm:"column:topic"`
	CommunityID string    `gorm:"foreign_key;column:community_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *CommunityPost) TableName() string {
	return "community_posts"
}
