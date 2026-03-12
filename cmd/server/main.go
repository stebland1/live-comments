package main

import (
	"log/slog"
	"os"

	"github.com/stebland1/live-comments/internal/config"
	httpapi "github.com/stebland1/live-comments/internal/transport/http"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	server := httpapi.NewServer(cfg)

	logger.Info("starting server", "host", cfg.Server.Host, "port", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("starting server", "host", cfg.Server.Host, "port", cfg.Server.Port)
		os.Exit(1)
	}
}
