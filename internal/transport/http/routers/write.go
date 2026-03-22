package routers

import (
	"net/http"

	"github.com/stebland1/live-comments/internal/transport/http/handlers"
)

func NewWriteRouter(commentHandler *handlers.CommentHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /comment", commentHandler.CreateComment)

	return mux
}
