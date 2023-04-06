package models

import "gorm.io/gorm"

type SiteChunk struct {
	gorm.Model
	Text      string
	URL       string
	TextIndex int
}
