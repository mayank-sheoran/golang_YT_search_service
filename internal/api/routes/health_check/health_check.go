package health_check

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handler(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
