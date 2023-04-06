package models

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Title      string
	VideoID    string `gorm:"unique"`
	Transcript string
}
