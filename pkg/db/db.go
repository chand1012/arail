package db

import (
	"os/user"
	"path/filepath"
	"strings"

	"github.com/chand1012/arail/pkg/db/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func New() (*Database, error) {
	db := &Database{}
	err := db.connect()
	if err != nil {
		return nil, err
	}
	err = db.DB.AutoMigrate(&models.SiteChunk{}, &models.Video{}, &models.Summary{}, &models.Chat{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewTemp() (*Database, error) {
	db := &Database{}
	err := db.connectMemory()
	if err != nil {
		return nil, err
	}
	err = db.DB.AutoMigrate(&models.SiteChunk{}, &models.Video{}, &models.Summary{}, &models.Chat{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *Database) connect() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	homeDir := currentUser.HomeDir
	dbPath := filepath.Join(homeDir, ".arail", "arail.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	d.DB = db

	if err != nil {
		return err
	}

	return nil
}

func (d *Database) connectMemory() error {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})

	d.DB = db

	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Close() error {
	db, err := d.DB.DB()
	if err != nil {
		return err
	}

	err = db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) PostSite(text models.SiteChunk) error {
	err := d.DB.Create(&text).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) SearchText(search string) ([]models.SiteChunk, error) {
	var texts []models.SiteChunk

	err := d.DB.Where("LOWER(text) LIKE ?", "%"+strings.ToLower(search)+"%").Find(&texts).Error
	if err != nil {
		return nil, err
	}

	return texts, nil
}

func (d *Database) GetTextByURL(url string) ([]models.SiteChunk, error) {
	var texts []models.SiteChunk

	err := d.DB.Where("url = ?", url).Find(&texts).Error
	if err != nil {
		return nil, err
	}

	return texts, nil
}

func (d *Database) PostVideo(video models.Video) error {
	err := d.DB.Create(&video).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetVideo(videoID string) (models.Video, error) {
	var video models.Video

	err := d.DB.Where("video_id = ?", videoID).First(&video).Error
	if err != nil {
		return models.Video{}, err
	}

	return video, nil
}

func (d *Database) GetSummaryByURL(url string) (models.Summary, error) {
	var summary models.Summary

	err := d.DB.Where("url = ?", url).First(&summary).Error
	if err != nil {
		return models.Summary{}, err
	}

	return summary, nil
}

func (d *Database) SearchSummary(q string) ([]models.Summary, error) {
	var summaries []models.Summary

	err := d.DB.Where("LOWER(summary) LIKE ?", "%"+strings.ToLower(q)+"%").Find(&summaries).Error
	if err != nil {
		return nil, err
	}

	return summaries, nil
}

// this is broken. Needs fixed
func (d *Database) SearchSummarySlice(qs []string) ([]models.Summary, error) {
	var summaries []models.Summary
	query := d.DB
	for _, q := range qs {
		query = query.Where("LOWER(summary) LIKE ?", "%"+strings.ToLower(q)+"%")
	}
	query = query.Order("created_at DESC")
	err := query.Find(&summaries).Error

	return summaries, err
}

func (d *Database) PostSummary(summary models.Summary) error {
	err := d.DB.Create(&summary).Error
	if err != nil {
		return err
	}

	return nil
}

// get the latest n chats with offset
func (d *Database) GetChats(n, offset int) ([]models.Chat, error) {
	var chats []models.Chat

	err := d.DB.Order("created_at DESC").Limit(n).Offset(offset).Find(&chats).Error
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (d *Database) GetChat(offset int) (models.Chat, error) {
	var chat models.Chat

	err := d.DB.Order("created_at DESC").Offset(offset).First(&chat).Error
	if err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

func (d *Database) PostChat(chat models.Chat) error {
	err := d.DB.Create(&chat).Error
	if err != nil {
		return err
	}

	return nil
}
