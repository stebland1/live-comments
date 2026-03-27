package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/stebland1/live-comments/internal/config"
	"github.com/stebland1/live-comments/internal/infra/redis"
	"github.com/stebland1/live-comments/internal/stream"
	httpapi "github.com/stebland1/live-comments/internal/transport/http"
	"github.com/stebland1/live-comments/internal/transport/http/handlers"
	"github.com/stebland1/live-comments/internal/transport/http/routers"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Load(); err != nil {
		logger.Info("no environment variable file", "err", err)
	}

	cfg := config.Load()

	commentSubscriber := redis.NewCommentSubscriber(cfg)
	streamService := stream.NewService(commentSubscriber)
	streamHandler := handlers.NewStreamHandler(streamService, logger)
	router := routers.NewStreamRouter(streamHandler)
	server := httpapi.NewStreamServer(cfg, router)

	logger.Info("starting server", "host", cfg.StreamServer.Host, "port", cfg.StreamServer.Port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error("failed to start server", "host", cfg.StreamServer.Host, "port", cfg.StreamServer.Port, "err", err)
		os.Exit(1)
	}
}
