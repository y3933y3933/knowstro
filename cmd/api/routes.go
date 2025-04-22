package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	r := gin.Default()
	{
		v1 := r.Group("/v1")
		v1.GET("/healthz", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version":     version,
				"status":      "available",
				"environment": app.config.env,
			})
		})
	}

	return r
}
