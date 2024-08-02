package models

import "gorm.io/gorm"

// Video model
type Video struct {
	gorm.Model
	Title     string
	URL       string
	Body      string
	UserID    uint // foreign key
	User      User
	Videosize int64
	MimeType  string
}
