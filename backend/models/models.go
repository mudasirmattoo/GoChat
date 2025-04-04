package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password"`
}

type Post struct {
	ID        uint           `gorm:"primaryKey" json:"post_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ImageURL  string         `json:"image_url,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Message struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SenderID   uint           `gorm:"not null" json:"sender_id"`
	ReceiverID uint           `gorm:"not null" json:"receiver_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type Like struct {
	gorm.Model
	PostID uint `gorm:"not null" json:"post_id"`
	UserID uint `gorm:"not null" json:"user_id"`
}

type Comments struct {
	gorm.Model
	PostID  uint   `gorm:"not null" json:"post_id"`
	Content string `gorm:"not null" json:"content"`
	UserID  uint   `gorm:"not null" json:"user_id"`
}

type PasswordReset struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
