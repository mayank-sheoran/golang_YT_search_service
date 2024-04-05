package yt_search

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/api/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/api/routes/yt_search/responses"
	"net/http"
)

func Handler(router *gin.Engine, ctx context.Context) {
	router.GET(
		"/videos/search", func(c *gin.Context) {
			searchQuery := c.Query("search-query")
			var response = responses.VideoSearchResponse{}
			response.Generate(searchQuery, ctx)
			c.JSON(
				http.StatusOK, models.NewGenericResponse(200, response, "success"),
			)
		},
	)
}
