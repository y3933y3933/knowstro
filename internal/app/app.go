package app

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/y3933y3933/knowstro/internal/api"
	"github.com/y3933y3933/knowstro/internal/mailer"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/migrations"
)

const version = "1.0.0"

type config struct {
	Port int
	Env  string
	SMTP smtpConfig
}

type smtpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type Application struct {
	Config              config
	Logger              *slog.Logger
	DB                  *sql.DB
	ResourceTypeHandler *api.ResourceTypeHandler
	UserHandler         *api.UserHandler
	Mailer              *mailer.Mailer
}

func NewApplication() (*Application, error) {
	cfg := loadConfig()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	mailer, err := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// handlers
	resourceTypeHandler := api.NewResourceTypeHandler(store.NewPostgresResourceTypeStore(pgDB), logger)
	userHandler := api.NewUserHandler(store.NewPostgresUserStore(pgDB), logger, mailer)

	app := &Application{
		Config:              cfg,
		Logger:              logger,
		DB:                  pgDB,
		ResourceTypeHandler: resourceTypeHandler,
		UserHandler:         userHandler,
		Mailer:              mailer,
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

func defaultString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func defaultInt(key string, fallback int) int {
	if s := os.Getenv(key); s != "" {
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}
	return fallback
}

func loadConfig() config {
	cfg := config{}
	fmt.Println("os Getenv", os.Getenv("SMTP_USERNAME"))
	flag.IntVar(&cfg.Port, "port", defaultInt("PORT", 8080), "API server port")
	flag.StringVar(&cfg.Env, "env", defaultString("ENV", "development"), "Environment (dev|prod)")
	flag.StringVar(&cfg.SMTP.Host, "smtp-host", defaultString("SMTP_HOST", "sandbox.smtp.mailtrap.io"), "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", defaultInt("SMTP_PORT", 25), "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", defaultString("SMTP_USERNAME", ""), "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", defaultString("SMTP_PASSWORD", ""), "SMTP password")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", defaultString("SMTP_SENDER", ""), "SMTP sender")
	flag.Parse()
	return cfg
}
