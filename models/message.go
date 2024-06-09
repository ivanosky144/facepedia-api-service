package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID        int       `gorm:"primary_key;column:id"`
	SenderID  int       `gorm:"foreign_key;column:sender_id"`
	Text      int       `gorm:"column:text"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Message) TableName() string {
	return "messages"
}
