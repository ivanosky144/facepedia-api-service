package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ID          int       `gorm:"primary_key;column:id"`
	UserID      int       `gorm:"foreign_key;column:user_id"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:desc"`
	Image       string    `gorm:"column:image"`
	Likes       []int     `gorm:"column:likes"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (p *Post) TableName() string {
	return "posts"
}
