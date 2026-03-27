package httpapi

import (
	"fmt"
	"net/http"

	"github.com/stebland1/live-comments/internal/config"
)

func NewWriteServer(cfg config.Config, handler *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.WriteServer.Host, cfg.WriteServer.Port),
		Handler: handler,
	}
}

func NewStreamServer(cfg config.Config, handler *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.StreamServer.Host, cfg.StreamServer.Port),
		Handler: handler,
	}
}
