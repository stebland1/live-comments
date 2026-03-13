package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebland1/live-comments/internal/config"
	"github.com/stebland1/live-comments/internal/repo"
	httpapi "github.com/stebland1/live-comments/internal/transport/http"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Info("no environment variable file", "err", err)
	}

	cfg := config.Load()

	_, err := repo.NewPostgresStore(cfg)
	if err != nil {
		logger.Error("creating db", "err", err)
		os.Exit(1)
	}

	server := httpapi.NewServer(cfg)

	logger.Info("starting server", "host", cfg.Server.Host, "port", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("starting server", "host", cfg.Server.Host, "port", cfg.Server.Port)
		os.Exit(1)
	}
}
