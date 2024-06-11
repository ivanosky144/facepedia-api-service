package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID  int
	PostID  int
	Content string
	Votes   []int
	User    User
	Post    Post
}

func (c *Comment) TableName() string {
	return "comments"
}
