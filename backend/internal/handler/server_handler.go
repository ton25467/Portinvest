package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"portinves/internal/domain"
	"portinves/internal/handler/middleware"
	"portinves/internal/service"
)

type ServerHandler struct {
	serverSrv *service.ServerService
}

func NewServerHandler(serverSrv *service.ServerService) *ServerHandler {
	return &ServerHandler{serverSrv: serverSrv}
}

type createServerReq struct {
	Name string            `json:"name"`
	Host string            `json:"host"`
	Port int               `json:"port"`
	Type domain.ServerType `json:"type"`
}

type createServiceCheckReq struct {
	Name            string            `json:"name"`
	Endpoint        string            `json:"endpoint"`
	Method          domain.HTTPMethod `json:"method"`
	ExpectedStatus  int               `json:"expected_status"`
	IntervalSeconds int               `json:"interval_seconds"`
}

func (h *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	var req createServerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	srv, err := h.serverSrv.CreateServer(r.Context(), userID, req.Name, req.Host, req.Port, req.Type)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, srv)
}

func (h *ServerHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	servers, err := h.serverSrv.ListServers(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, servers)
}

func (h *ServerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing server ID"})
		return
	}

	srv, err := h.serverSrv.GetServer(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, srv)
}

func (h *ServerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing server ID"})
		return
	}

	var req createServerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	srv, err := h.serverSrv.GetServer(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	srv.Name = req.Name
	srv.Host = req.Host
	srv.Port = req.Port
	srv.Type = req.Type

	if err := h.serverSrv.UpdateServer(r.Context(), srv); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, srv)
}

func (h *ServerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing server ID"})
		return
	}

	if err := h.serverSrv.DeleteServer(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func (h *ServerHandler) CreateCheck(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	if serverID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing server ID"})
		return
	}

	var req createServiceCheckReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	check, err := h.serverSrv.CreateServiceCheck(
		r.Context(),
		serverID,
		req.Name,
		req.Endpoint,
		req.Method,
		req.ExpectedStatus,
		req.IntervalSeconds,
	)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, check)
}

func (h *ServerHandler) ListChecks(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	if serverID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing server ID"})
		return
	}

	checks, err := h.serverSrv.ListServiceChecks(r.Context(), serverID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, checks)
}

func (h *ServerHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	serverID := chi.URLParam(r, "id")
	if serverID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing server ID"})
		return
	}

	limit := 50
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
			limit = val
		}
	}

	logs, err := h.serverSrv.GetUptimeLogs(r.Context(), serverID, limit)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, logs)
}
