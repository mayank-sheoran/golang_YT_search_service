package syncronizations

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"strings"
	"time"
)

func getLatestPublishedAtInVideosMetaDataIndex(client *elasticsearch.Client, ctx context.Context) time.Time {
	var defaultTime time.Time
	defaultTime, _ = time.Parse("2006-01-02", "1900-01-01")
	query := `{
        "size": 1,
        "sort": [{ "publishedAt": { "order": "desc" }}],
        "_source": ["publishedAt"]
    }`
	res, err := client.Search(
		client.Search.WithContext(ctx),
		client.Search.WithIndex("videos_meta_data"),
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

func (ess *EsSyncronizations) SyncVideosMetaDataIndex(client *elasticsearch.Client, ctx context.Context) {
	for {
		latestPublishedAtPresent := getLatestPublishedAtInVideosMetaDataIndex(client, ctx)
		var videosMetaData []models.VideoMetaData
		result := db.YtSearchServiceDb.Find(&videosMetaData, "published_at > ?", latestPublishedAtPresent)
		log.HandleError(result.Error, ctx, false)

		var buf bytes.Buffer
		if len(videosMetaData) == 0 {
			time.Sleep(time.Second * 10)
			continue
		}

		for _, data := range videosMetaData {
			meta := map[string]interface{}{
				"index": map[string]interface{}{
					"_index": "videos_meta_data",
					"_id":    data.ID,
				},
			}
			metaStr, err := json.Marshal(meta)
			log.HandleError(err, ctx, false)
			dataStr, err := json.Marshal(data)
			log.HandleError(err, ctx, false)
			buf.Grow(len(metaStr) + len(dataStr) + 2) // for newline characters
			buf.Write(metaStr)
			buf.WriteByte('\n')
			buf.Write(dataStr)
			buf.WriteByte('\n')
		}
		res, err := client.Bulk(
			&buf, client.Bulk.WithIndex("videos_meta_data"),
		)
		log.HandleError(err, ctx, false)
		defer res.Body.Close()
		if res.IsError() {
			log.Error("Error with bulk request", ctx)
		} else {
			log.Info("Bulk push completed successfully for videos_meta_data | Elastic search", ctx)
		}
		time.Sleep(time.Second * 10)
	}
}
