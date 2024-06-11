package models

import (
	"gorm.io/gorm"
)

type Moderator struct {
	gorm.Model
	CommunityID       int
	CommunityMemberID int
	Community         Community
	CommunityMember   CommunityMember
}

func (m *Moderator) TableName() string {
	return "moderators"
}
