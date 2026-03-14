package httpapi

import (
	"fmt"
	"net/http"

	"github.com/stebland1/live-comments/internal/config"
	"github.com/stebland1/live-comments/internal/transport/http/handlers"
)

func NewServer(cfg config.Config, commentHandler *handlers.CommentHandler) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: NewRouter(commentHandler),
	}
}
