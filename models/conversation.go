package models

import (
	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	Title     string
	CreatorID int
	Creator   User
}

func (c *Conversation) TableName() string {
	return "conversations"
}
