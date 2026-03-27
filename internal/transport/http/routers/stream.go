package routers

import (
	"net/http"

	"github.com/stebland1/live-comments/internal/transport/http/handlers"
)

func NewStreamRouter(streamHandler *handlers.StreamHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /video/{videoID}/stream", streamHandler.Stream)

	return mux
}
