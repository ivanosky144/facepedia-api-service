package model

import (
	"time"

	"gorm.io/gorm"
)

type CommunityPost struct {
	gorm.Model
	ID          int       `gorm:"primary_key;column:id"`
	PostID      string    `gorm:"primary_key;column:post_id"`
	CommunityID string    `gorm:"primary_key;column:community_id"`
	Mark        string    `gorm:"primary_key;column:mark"`
	Topic       string    `gorm:"primary_key;column:topic"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *CommunityPost) TableName() string {
	return "community_posts"
}
