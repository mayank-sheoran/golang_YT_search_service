package db

import (
	"context"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
)

func gormAutoMigrations(ctx context.Context) {
	err := YtSearchServiceDb.AutoMigrate(&models.VideoMetaData{})
	log.HandleError(err, ctx, true)
}
