package repository

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"strings"
	"time"
)

type VideoMetaDataIndexRepository struct {
	IndexName string
}

var (
	VideoMetaDataIndexRepoClient = &VideoMetaDataIndexRepository{
		IndexName: "videos_meta_data",
	}
)

func (vmdir *VideoMetaDataIndexRepository) SearchInTitleAndDescription(
	esClient *elasticsearch.Client, searchQuery string, ctx context.Context,
) []models.VideoMetaData {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  searchQuery,
				"fields": []string{"title", "description"},
			},
		},
	}
	queryJSON, err := json.Marshal(query)
	log.HandleError(err, ctx, false)

	res, err := esClient.Search(
		esClient.Search.WithContext(context.Background()),
		esClient.Search.WithIndex(vmdir.IndexName), // Change this to your Elasticsearch index name
		esClient.Search.WithBody(strings.NewReader(string(queryJSON))),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
	)
	log.HandleError(err, ctx, false)
	defer res.Body.Close()

	var searchResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		log.HandleError(err, ctx, false)
	}

	hits := searchResponse["hits"].(map[string]interface{})["hits"].([]interface{})
	var videosMetaData []models.VideoMetaData
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		publishedAtStr := source["publishedAt"].(string)
		publishedAt, err := time.Parse(time.RFC3339, publishedAtStr)
		log.HandleError(err, ctx, false)
		videoMetaData := models.VideoMetaData{
			ID:           source["ID"].(string),
			Title:        source["title"].(string),
			Description:  source["description"].(string),
			PublishedAt:  publishedAt,
			ThumbnailURL: source["thumbnailUrl"].(string),
		}
		videosMetaData = append(videosMetaData, videoMetaData)
	}
	return videosMetaData
}

func (vmdir *VideoMetaDataIndexRepository) GetLatestPublishedAtInVideosMetaDataIndex(
	client *elasticsearch.Client, ctx context.Context,
) time.Time {
	var defaultTime time.Time
	defaultTime, _ = time.Parse("2006-01-02", "1900-01-01")
	query := `{
        "size": 1,
        "sort": [{ "publishedAt": { "order": "desc" }}],
        "_source": ["publishedAt"]
    }`
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex(vmdir.IndexName),
		client.Search.WithBody(strings.NewReader(query)),
	)
	log.HandleError(err, ctx, true)
	defer res.Body.Close()

	if res.IsError() {
		return defaultTime
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return defaultTime
	}
	hits := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) > 0 {
		firstHit := hits[0].(map[string]interface{})
		source := firstHit["_source"].(map[string]interface{})
		publishedAtStr := source["publishedAt"].(string)
		publishedAt, err := time.Parse(time.RFC3339, publishedAtStr)
		if err != nil {
			return defaultTime
		}
		return publishedAt
	}
	return defaultTime
}
