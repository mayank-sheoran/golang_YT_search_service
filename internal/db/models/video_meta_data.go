package models

import (
	"time"
)

type VideoMetaData struct {
	ID           uint64 `gorm:"primarykey"`
	Title        string
	Description  string
	PublishedAt  time.Time
	ThumbnailURL string
}
