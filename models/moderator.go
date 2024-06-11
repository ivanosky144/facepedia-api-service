package models

import (
	"time"

	"gorm.io/gorm"
)

type Moderator struct {
	gorm.Model
	ID                int       `gorm:"primary_key;column:id"`
	CommunityID       int       `gorm:"column:community_id"`
	CommunityMemberID int       `gorm:"column:communitymember_id"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (m *Moderator) TableName() string {
	return "moderators"
}
