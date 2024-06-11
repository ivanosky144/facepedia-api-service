package models

import (
	"gorm.io/gorm"
)

type CommunityMember struct {
	gorm.Model
	UserID      User
	CommunityID Community
	User        int
	Community   int
}

func (c *CommunityMember) TableName() string {
	return "community_members"
}
