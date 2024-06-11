package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID        int
	Title         string
	Description   string
	Image         string
	Likes         []int
	User          User
	Comment       []Comment
	CommunityPost []CommunityPost
}

func (p *Post) TableName() string {
	return "posts"
}
