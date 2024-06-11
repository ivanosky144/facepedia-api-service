package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID int
	Text     int
	Sender   User
}

func (c *Message) TableName() string {
	return "messages"
}
