package models

import (
	"gorm.io/gorm/clause"
	"time"
)

type VideoMetaData struct {
	ID           string `gorm:"primarykey"`
	Title        string
	Description  string
	PublishedAt  time.Time
	ThumbnailURL string
}

func (VideoMetaData) TableName() string {
	return "video_meta_data"
}

func VideoMetaDataUpsertClause() clause.Expression {
	return clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns(
			[]string{
				"title", "description", "published_at", "thumbnail_url",
			},
		),
	}
}
