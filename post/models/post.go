package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"post_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ImageURL  string         `json:"image_url,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
