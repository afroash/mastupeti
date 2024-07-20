package models

import "gorm.io/gorm"

// User model
type User struct {
	gorm.Model
	Email    string `gorm:"unique"`
	Password string
	Username string `gorm:"unique"`
}
