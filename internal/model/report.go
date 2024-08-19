package model

import (
	"time"

	"gorm.io/gorm"
)

type Report struct {
	gorm.Model
	ID        int       `gorm:"primary_key;column:id"`
	PostID    *int      `gorm:"column:post_id;null"`
	CommentID *int      `gorm:"column:comment_id;null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdateAt  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (r *Report) TableName() string {
	return "reports"
}
