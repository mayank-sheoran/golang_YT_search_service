package models

import (
	"gorm.io/gorm/clause"
	"time"
)

type VideoMetaData struct {
	ID           string    `gorm:"primarykey" json:"ID"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PublishedAt  time.Time `json:"publishedAt"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	CreatedAt    time.Time `json:"createdAt"`
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
