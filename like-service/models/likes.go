package models

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	PostID uint `gorm:"not null" json:"post_id"`
	UserID uint `gorm:"not null" json:"user_id"`
}
