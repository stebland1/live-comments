package main

import (
	"fmt"
	"log"

	"github.com/stebland1/live-comments/internal/config"
	httpapi "github.com/stebland1/live-comments/internal/transport/http"
)

func main() {
	cfg := config.Load()
	server := httpapi.NewServer(cfg)

	fmt.Printf("starting server at %s:%s\n", cfg.Server.Host, cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start server")
	}
}
