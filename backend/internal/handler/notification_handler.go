package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"portinves/internal/handler/middleware"
	"portinves/internal/service"
)

type NotificationHandler struct {
	notifSrv *service.NotificationService
}

func NewNotificationHandler(notifSrv *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSrv: notifSrv}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	unreadOnly := r.URL.Query().Get("unread_only") == "true"

	notifs, err := h.notifSrv.ListUserNotifications(r.Context(), userID, unreadOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(notifs)
}

func (h *NotificationHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	notifID := chi.URLParam(r, "id")
	if notifID == "" {
		http.Error(w, "missing notification id", http.StatusBadRequest)
		return
	}

	if err := h.notifSrv.MarkAsRead(r.Context(), notifID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
