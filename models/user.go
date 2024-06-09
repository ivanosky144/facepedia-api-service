package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             int       `gorm:"primary_key;column:id"`
	Username       string    `gorm:"column:username"`
	Email          string    `gorm:"column:email"`
	Password       string    `gorm:"column:password"`
	ProfilePicture string    `gorm:"column:profile_picture"`
	CoverPicture   string    `gorm:"column:cover_picture"`
	Followers      []string  `gorm:"column:followers"`
	Followings     []string  `gorm:"column:followings"`
	SocialPoint    int       `gorm:"column:social_point"`
	Desc           string    `gorm:"column:desc"`
	City           string    `gorm:"column:city"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
