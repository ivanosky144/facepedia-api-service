package models

import (
	"gorm.io/gorm"
)

type Participant struct {
	gorm.Model
	ConversationID int
	UserID         int
	Conversation   Conversation
	User           User
}

func (c *Participant) TableName() string {
	return "participants"
}
