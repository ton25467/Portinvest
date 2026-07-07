package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"portinves/internal/handler/middleware"
	"portinves/internal/service"
)

type PortfolioHandler struct {
	portfolioSrv *service.PortfolioService
}

func NewPortfolioHandler(portfolioSrv *service.PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{portfolioSrv: portfolioSrv}
}

type createPortfolioReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Currency    string `json:"currency"`
}

func (h *PortfolioHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	var req createPortfolioReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	p, err := h.portfolioSrv.CreatePortfolio(r.Context(), userID, req.Name, req.Description, req.Currency)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

func (h *PortfolioHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		return
	}

	portfolios, err := h.portfolioSrv.ListPortfolios(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, portfolios)
}

func (h *PortfolioHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing portfolio ID"})
		return
	}

	p, err := h.portfolioSrv.GetPortfolio(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

func (h *PortfolioHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing portfolio ID"})
		return
	}

	summary, err := h.portfolioSrv.GetPortfolioSummary(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, summary)
}
