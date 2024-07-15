package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Title string
	URL   string
	Body  string
}
