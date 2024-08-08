package model

import (
	"time"

	"gorm.io/gorm"
)

type UserFollow struct {
	gorm.Model
	ID          int       `gorm:"primary_key;column:id"`
	FollowerID  int       `gorm:"column:follower_id"`
	FollowingID int       `gorm:"column:following_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (uf *UserFollow) TableName() string {
	return "user_follows"
}
