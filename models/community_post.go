package models

import (
	"gorm.io/gorm"
)

type CommunityPost struct {
	gorm.Model
	PostID      string
	CommunityID string
	Mark        string
	Topic       string
	Post        Post
	Community   Community
}

func (c *CommunityPost) TableName() string {
	return "community_posts"
}
