package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/app"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.GET("/healthz", app.HealthCheck)

	return r
}
