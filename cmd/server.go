package cmd

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/api"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/common/enums/env"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils/log"
	"os"
)

func StartServer(ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(corsMiddleware())

	// Add Routes
	api.HandleRoutes(router, ctx)

	// Start server
	log.Info("server started on port: "+os.Getenv(env.Port), ctx)
	err := router.Run(":" + os.Getenv(env.Port))
	log.HandleError(err, ctx, true)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := []string{"http://localhost:3000"}
		origin := c.GetHeader("Origin")
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}
		c.Next()
	}
}
