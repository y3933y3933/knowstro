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

	r.Use(app.UserMiddleware.Authenticate())

	{
		v1 := r.Group("v1")
		v1.GET("/healthz", app.HealthCheck)
		{
			{
				types := v1.Group("/types")
				types.GET("/:id", app.ResourceTypeHandler.GetTypeByID)
				types.POST("", app.ResourceTypeHandler.CreateType)
				types.PUT("/:id", app.ResourceTypeHandler.UpdateType)
				types.DELETE("/:id", app.ResourceTypeHandler.DeleteType)
				types.GET("", app.ResourceTypeHandler.ListTypes)
				types.DELETE("/reset", app.ResourceTypeHandler.ResetTypes)
			}

			{
				users := v1.Group("/users")
				users.POST("", app.UserHandler.HandleRegisterUser)
				users.PUT("/activated", app.UserHandler.HandlerActivateUser)
			}

			{
				tokens := v1.Group("tokens")
				tokens.POST("/authentication", app.TokenHandler.HandleCreateToken)
			}

		}

	}

	return r
}
