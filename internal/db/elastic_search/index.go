package elastic_search

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/elastic_search/syncronizations"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
)

type ElasticSearch struct{}

var (
	ElasticClient *elasticsearch.Client
)

func (es *ElasticSearch) Init(ctx context.Context) {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	log.HandleError(err, ctx, true)
	ElasticClient = esClient
	es.syncWithPostgresDB(ctx)
}

func (es *ElasticSearch) syncWithPostgresDB(ctx context.Context) {
	esSync := syncronizations.EsSyncronizations{}
	go esSync.SyncVideosMetaDataIndex(ElasticClient, ctx)
}
