package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/y3933y3933/knowstro/internal/api"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/migrations"
)

const version = "1.0.0"

type config struct {
	Port int
	Env  string
}

type Application struct {
	Config              config
	Logger              *slog.Logger
	DB                  *sql.DB
	ResourceTypeHandler *api.ResourceTypeHandler
}

func NewApplication() (*Application, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// handlers
	resourceTypeHandler := api.NewResourceTypeHandler(store.NewResourceTypeStore(pgDB))

	app := &Application{
		Logger:              logger,
		DB:                  pgDB,
		ResourceTypeHandler: resourceTypeHandler,
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
