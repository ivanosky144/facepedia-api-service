package model

import (
	"time"

	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	ID           int           `gorm:"primary_key;column:id"`
	Title        string        `gorm:"column:title"`
	UserID       int           `gorm:"column:user_id"`
	Participants []Participant `gorm:"foreignKey:ConversationID"`
	CreatedAt    time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Conversation) TableName() string {
	return "conversations"
}
