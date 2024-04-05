package yt_search

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/api/models"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/api/routes/yt_search/responses"
	"net/http"
	"strconv"
)

func Handler(router *gin.Engine, ctx context.Context) {
	router.GET(
		"/videos/search", func(c *gin.Context) {
			searchQuery := c.Query("search-query")
			pageNumber, err := strconv.Atoi(c.Query("page-number"))
			if err != nil {
				pageNumber = 1
			}
			perPageLimit, err := strconv.Atoi(c.Query("page-limit"))
			if err != nil {
				perPageLimit = 100
			}
			var response = responses.VideoSearchResponse{}
			response.Generate(searchQuery, pageNumber, perPageLimit, ctx)
			c.JSON(
				http.StatusOK, models.NewGenericResponse(200, response, "success"),
			)
		},
	)
}
