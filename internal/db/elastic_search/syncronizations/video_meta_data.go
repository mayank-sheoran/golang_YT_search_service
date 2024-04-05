package syncronizations

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/elastic_search/repository"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"time"
)

func (ess *EsSyncronizations) SyncVideosMetaDataIndex(client *elasticsearch.Client, ctx context.Context) {
	for {
		latestPublishedAtPresent := repository.VideoMetaDataIndexRepoClient.GetLatestPublishedAtInVideosMetaDataIndex(
			client, ctx,
		)
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
