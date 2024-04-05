package responses

import (
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/elastic_search"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/elastic_search/repository"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"golang.org/x/net/context"
)

type VideoSearchResponse struct {
	Videos []models.VideoMetaData
}

func (vsr *VideoSearchResponse) Generate(searchQuery string, page, perPage int, ctx context.Context) {
	if searchQuery == "" {
		var videosMetaData []models.VideoMetaData
		result := db.YtSearchServiceDb.Model(&videosMetaData).
			Order("created_at desc").
			Offset((page - 1) * perPage).
			Limit(perPage).
			Find(&videosMetaData)
		log.HandleError(result.Error, ctx, false)
		vsr.Videos = videosMetaData
		return
	}
	vsr.Videos = repository.VideoMetaDataIndexRepoClient.SearchInTitleAndDescription(
		elastic_search.ElasticClient, searchQuery, page, perPage, ctx,
	)
}
