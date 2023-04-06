package models

import "gorm.io/gorm"

type Summary struct {
	gorm.Model
	URL     string `gorm:"unique"`
	Summary string
}
