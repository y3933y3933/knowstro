package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/app"
)

func SetupRoutes(app *app.Application) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20)
		c.Next()
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/healthz", app.HealthCheck)

		types := v1.Group("/types")
		{
			types.GET("/:id", app.ResourceTypeHandler.GetTypeByID)
			types.POST("", app.ResourceTypeHandler.CreateType)
			types.PUT("/:id", app.ResourceTypeHandler.UpdateType)
			types.DELETE("/:id", app.ResourceTypeHandler.DeleteType)
			types.GET("", app.ResourceTypeHandler.ListTypes)
		}

	}

	return r
}
