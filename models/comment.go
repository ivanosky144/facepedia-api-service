package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        int       `gorm:"primary_key;column:id"`
	UserID    int       `gorm:"column:user_id"`
	PostID    int       `gorm:"column:post_id"`
	ParentID  int       `gorm:"column:parent_id"`
	Content   string    `gorm:"column:content"`
	Replies   []Comment `gorm:"foreignKey:ParentID;references:ID"`
	Parent    *Comment  `gorm:"foreignKey:ParentID;references:ID"`
	Votes     []*User   `gorm:"many2many:user_votes;"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Comment) TableName() string {
	return "comments"
}
