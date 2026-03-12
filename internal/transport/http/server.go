package httpapi

import (
	"fmt"
	"net/http"

	"github.com/stebland1/live-comments/internal/config"
)

func NewServer(cfg config.Config) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: NewRouter(),
	}
}
