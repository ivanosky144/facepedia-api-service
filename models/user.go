package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         string
	Email            string
	Password         string
	ProfilePicture   string
	CoverPicture     string
	Followers        []string
	Followings       []string
	SocialPoint      int
	Desc             string
	Country          string
	Posts            []Post
	Comments         []Comment
	CommunityMembers []CommunityMember
	Conversations    []Conversation
	Participants     []Participant
}

func (u *User) TableName() string {
	return "users"
}
