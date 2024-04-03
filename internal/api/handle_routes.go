package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/api/routes/health_check"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
)

func HandleRoutes(router *gin.Engine, ctx context.Context) {
	ctx = utils.GetContextWithFlowName(ctx, "routes setup")
	health_check.Handler(router)
	log.Info("all routes are active.", ctx)
}
