package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const version = "1.0.0"

type config struct {
	Port int
	Env  string
}

type Application struct {
	Config config
	Logger *slog.Logger
}

func NewApplication() (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &Application{
		Logger: logger,
	}

	return app, nil
}

func (a *Application) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     version,
		"status":      "available",
		"environment": a.Config.Env,
	})
}
