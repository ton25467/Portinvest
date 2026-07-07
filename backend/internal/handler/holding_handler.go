package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"portinves/internal/domain"
	"portinves/internal/service"
)

type HoldingHandler struct {
	portfolioSrv *service.PortfolioService
}

func NewHoldingHandler(portfolioSrv *service.PortfolioService) *HoldingHandler {
	return &HoldingHandler{portfolioSrv: portfolioSrv}
}

type addTransactionReq struct {
	Symbol     string                 `json:"symbol"`
	Name       string                 `json:"name"`
	AssetType  domain.AssetType       `json:"asset_type"`
	Type       domain.TransactionType `json:"type"`
	Quantity   float64                `json:"quantity"`
	Price      float64                `json:"price"`
	Fee        float64                `json:"fee"`
	Notes      string                 `json:"notes"`
	ExecutedAt string                 `json:"executed_at"` // ISO8601 string
}

func (h *HoldingHandler) ListHoldings(w http.ResponseWriter, r *http.Request) {
	portfolioID := chi.URLParam(r, "id")
	if portfolioID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing portfolio ID"})
		return
	}

	holdings, err := h.portfolioSrv.GetHoldings(r.Context(), portfolioID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, holdings)
}

func (h *HoldingHandler) AddTransaction(w http.ResponseWriter, r *http.Request) {
	portfolioID := chi.URLParam(r, "id")
	if portfolioID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "missing portfolio ID"})
		return
	}

	var req addTransactionReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	var executedAt time.Time
	if req.ExecutedAt != "" {
		var err error
		executedAt, err = time.Parse(time.RFC3339, req.ExecutedAt)
		if err != nil {
			// Try parsing as simple date
			executedAt, err = time.Parse("2006-01-02", req.ExecutedAt)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid executed_at format, must be RFC3339 or YYYY-MM-DD"})
				return
			}
		}
	} else {
		executedAt = time.Now()
	}

	tx, err := h.portfolioSrv.AddTransaction(
		r.Context(),
		portfolioID,
		req.Symbol,
		req.Name,
		req.AssetType,
		req.Type,
		req.Quantity,
		req.Price,
		req.Fee,
		req.Notes,
		executedAt,
	)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, tx)
}
