package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/store"
)

const version = "1.0.0"

type config struct {
	Port int
	Env  string
}

type Application struct {
	Config config
	Logger *slog.Logger
	DB     *sql.DB
}

func NewApplication() (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	app := &Application{
		Logger: logger,
		DB:     pgDB,
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
