package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type StreamService interface {
	Subscribe(ctx context.Context, videoID int64) <-chan string
}

type StreamHandler struct {
	service StreamService
	logger  *slog.Logger
}

func NewStreamHandler(service StreamService, logger *slog.Logger) *StreamHandler {
	return &StreamHandler{service: service, logger: logger}
}

func (h *StreamHandler) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.logger.Error("response writer doesn't support flush")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()
	videoIDStr := r.PathValue("videoID")
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		h.logger.Error("cannot convert 'videoID' into int64 type")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	ch := h.service.Subscribe(ctx, videoID)

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Connection", "Keep-Alive")
	w.Header().Add("Cache-Control", "no-cache")

	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		case <-ctx.Done():
			return
		}
	}
}
