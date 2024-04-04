package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mayank-sheoran/golang_YT_search_service/cmd"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/db"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/service"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
)

func loadEnvFile(ctx context.Context) {
	err := godotenv.Load()
	log.HandleError(err, ctx, true)
}

func main() {
	var ctx = utils.GetContextWithFlowName(context.Background(), "main")
	loadEnvFile(ctx)

	// connect database and elastic data store
	db.ConnectDatabase(utils.GetContextWithFlowName(ctx, "database connection"))
	db.ConnectElasticSearch(utils.GetContextWithFlowName(ctx, "elastic database connection"))

	// start HTTP server
	go cmd.StartServer(ctx)

	// start YT search service
	go service.NewYoutubeDataV3Client().Run()

	// To keep the service running
	select {}
}
