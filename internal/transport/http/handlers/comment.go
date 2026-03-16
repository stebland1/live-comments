package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID int64, videoID int64, content string) (int64, error)
}

type CommentHandler struct {
	service CommentService
	logger  *slog.Logger
}

func NewCommentHandler(service CommentService, logger *slog.Logger) *CommentHandler {
	return &CommentHandler{service: service, logger: logger}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		VideoID int64  `json:"video_id"`
		Content string `json:"content"`
	}

	var createCommentReq Request
	if err := json.NewDecoder(r.Body).Decode(&createCommentReq); err != nil {
		h.logger.Warn("invalid create comment json", "err", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		h.logger.Warn("invalid user id header", "err", err)
		http.Error(w, "missing user id in request", http.StatusBadRequest)
		return
	}

	commentID, err := h.service.CreateComment(
		r.Context(),
		userID,
		createCommentReq.VideoID,
		createCommentReq.Content,
	)

	if err != nil {
		h.logger.Error("creating comment", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	RespondJSON(w, http.StatusCreated, map[string]int64{
		"comment_id": commentID,
	})
}
