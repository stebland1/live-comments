package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebland1/live-comments/internal/comment"
	"github.com/stebland1/live-comments/internal/config"
	"github.com/stebland1/live-comments/internal/infra/postgres"
	httpapi "github.com/stebland1/live-comments/internal/transport/http"
	"github.com/stebland1/live-comments/internal/transport/http/handlers"

	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Info("no environment variable file", "err", err)
	}

	cfg := config.Load()

	commentRepo, err := postgres.NewCommentRepo(cfg)
	if err != nil {
		logger.Error("creating db", "err", err)
		os.Exit(1)
	}

	commentService := comment.NewService(commentRepo)
	commentHandler := handlers.NewCommentHandler(commentService)
	commentHandler := handlers.NewCommentHandler(commentService, logger)
	server := httpapi.NewServer(cfg, commentHandler)

	logger.Info("starting server", "host", cfg.Server.Host, "port", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("starting server", "host", cfg.Server.Host, "port", cfg.Server.Port, "err", err)
		os.Exit(1)
	}
}
