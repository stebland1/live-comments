package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebland1/live-comments/internal/comment"
	"github.com/stebland1/live-comments/internal/config"
	"github.com/stebland1/live-comments/internal/infra/postgres"
	"github.com/stebland1/live-comments/internal/infra/redis"
	httpapi "github.com/stebland1/live-comments/internal/transport/http"
	"github.com/stebland1/live-comments/internal/transport/http/handlers"
	"github.com/stebland1/live-comments/internal/transport/http/routers"

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

	commentPublisher := redis.NewCommentPublisher(cfg)
	commentService := comment.NewService(commentRepo, commentPublisher)
	commentHandler := handlers.NewCommentHandler(commentService, logger)
	router := routers.NewWriteRouter(commentHandler)
	server := httpapi.NewWriteServer(cfg, router)

	logger.Info("starting server", "host", cfg.WriteServer.Host, "port", cfg.WriteServer.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("starting server", "host", cfg.WriteServer.Host, "port", cfg.WriteServer.Port, "err", err)
		os.Exit(1)
	}
}
