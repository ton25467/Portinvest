package handler

import (
	"net/http"

	"portinves/internal/domain"
	"portinves/internal/handler/middleware"
)

type DashboardHandler struct {
	dashboardRepo domain.DashboardRepository
}

func NewDashboardHandler(dashboardRepo domain.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{dashboardRepo: dashboardRepo}
}

func (h *DashboardHandler) GetOverview(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	overview, err := h.dashboardRepo.GetOverview(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, overview)
}
