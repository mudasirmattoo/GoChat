package models

import "gorm.io/gorm"

type Comments struct {
	gorm.Model
	PostID  uint   `gorm:"not null" json:"post_id"`
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `gorm:"not null" json:"user_id"`
}
