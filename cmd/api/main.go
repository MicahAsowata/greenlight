package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spobly/greenlight/internal/config"
	"github.com/spobly/greenlight/internal/data"
	"github.com/spobly/greenlight/internal/mailer"
)

const version = "1.0.0"

type application struct {
	config config.Config
	logger *slog.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	cfg := &config.Config{}
	cfg.Parse()

	db, err := data.OpenDB(*cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	app := &application{
		config: *cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.SMTP.Port, cfg.SMTP.Host, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
	}
}
