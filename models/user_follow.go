package models

import (
	"time"

	"gorm.io/gorm"
)

type UserFollow struct {
	gorm.Model
	ID        int       `gorm:"primary_key;column:id"`
	UserID    int       `gorm:"column:user_id"`
	FollowID  int       `gorm:"column:follow_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (uf *UserFollow) TableName() string {
	return "user_follows"
}
