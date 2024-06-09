package models

import (
	"time"

	"gorm.io/gorm"
)

type CommunityMember struct {
	gorm.Model
	ID          int       `gorm:"primary_key;column:id"`
	UserID      int       `gorm:"foreign_key;column:user_id"`
	CommunityID int       `gorm:"foreign_key;column:community_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *CommunityMember) TableName() string {
	return "community_members"
}
