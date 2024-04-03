package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/mayank-sheoran/golang_YT_search_service/cmd"
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

	// start HTTP server
	go cmd.StartServer(ctx)

	// To keep the service running
	select {}
}
