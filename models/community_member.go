package models

import (
	"time"

	"gorm.io/gorm"
)

type CommunityMember struct {
	gorm.Model
	UserID      int       `gorm:"primary_key;column:user_id"`
	CommunityID int       `gorm:"primary_key;column:community_id"`
	IsModerator bool      `gorm:"column:is_moderator"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *CommunityMember) TableName() string {
	return "community_members"
}
