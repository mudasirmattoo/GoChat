package models

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Adds ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string `gorm:"unique;not null" json:"username"`
	Email      string `gorm:"unique;not null" json:"email"`
	Password   string `json:"password"`
}

//The json tags allow automatic JSON parsing in request/response handling.
