package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mayank-sheoran/golang_YT_search_service/cmd"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db/elastic_search"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/service"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"time"
)

func loadEnvFile(ctx context.Context) {
	err := godotenv.Load()
	log.HandleError(err, ctx, true)
}

func setServerTimeToIst(ctx context.Context) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	log.HandleError(err, ctx, true)
	time.Local = loc
}

func main() {
	var ctx = utils.GetContextWithFlowName(context.Background(), "main")
	loadEnvFile(ctx)
	setServerTimeToIst(ctx)

	// connect database and elastic data store
	db.ConnectDatabase(utils.GetContextWithFlowName(ctx, "database connection"))
	es := elastic_search.ElasticSearch{}
	es.Init(utils.GetContextWithFlowName(ctx, "elastic database connection"))

	// start HTTP server
	go cmd.StartServer(ctx)

	// start YT search service
	go service.NewYoutubeDataV3Client().Run()

	// To keep the service running
	select {}
}
