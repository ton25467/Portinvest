package handler

import (
	"log/slog"
	"net/http"

	"github.com/coder/websocket"

	webSocketHub "portinves/internal/websocket"
	"portinves/internal/service"
)

type WSHandler struct {
	hub     *webSocketHub.Hub
	authSrv *service.AuthService
}

func NewWSHandler(hub *webSocketHub.Hub, authSrv *service.AuthService) *WSHandler {
	return &WSHandler{hub: hub, authSrv: authSrv}
}

func (h *WSHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.URL.Query().Get("token")
	if tokenStr == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	claims, err := h.authSrv.ValidateToken(tokenStr)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow CORS in development
	})
	if err != nil {
		slog.Error("Failed to accept websocket connection", "error", err)
		return
	}

	client := webSocketHub.NewClient(h.hub, conn, claims.UserID)
	h.hub.Register(client)

	// Start read/write loops
	go client.WritePump(r.Context())
	go client.ReadPump(r.Context())
}
