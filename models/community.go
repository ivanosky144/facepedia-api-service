package models

import (
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name             string
	Description      string
	Rules            []string
	MembersCount     int
	LogoPicture      string
	CommunityMembers []CommunityMember
	Moderators       []Moderator
	CommunityPosts   []CommunityPost
}

func (c *Community) TableName() string {
	return "communities"
}
