package models

import (
	"time"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	ID           int       `gorm:"primary_key;column:id"`
	Name         string    `gorm:"column:name"`
	Description  string    `gorm:"column:desc"`
	Rules        []string  `gorm:"column:rules"`
	Members      []int     `gorm:"column:members"`
	MembersCount int       `gorm:column:members_count`
	Moderators   []int     `gorm:"column:moderators"`
	LogoPicture  string    `gorm:"column:logo_picture"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Community) TableName() string {
	return "communities"
}
