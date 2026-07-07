package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"portinves/internal/domain"
	"portinves/internal/handler/middleware"
	"portinves/internal/service"
)

type CashflowHandler struct {
	cashflowSrv *service.CashflowService
	syncSrv     *service.WebScraperSyncService
}

func NewCashflowHandler(cashflowSrv *service.CashflowService, syncSrv *service.WebScraperSyncService) *CashflowHandler {
	return &CashflowHandler{cashflowSrv: cashflowSrv, syncSrv: syncSrv}
}

type CreateCashflowReq struct {
	PortfolioID *string             `json:"portfolio_id,omitempty"`
	Type        domain.CashflowType `json:"type"`
	Amount      float64             `json:"amount"`
	Currency    string              `json:"currency"`
	Description string              `json:"description"`
}

func (h *CashflowHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateCashflowReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf, err := h.cashflowSrv.CreateCashflow(r.Context(), userID, req.PortfolioID, req.Type, req.Amount, req.Currency, req.Description, time.Now().UTC())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cf)
}

func (h *CashflowHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	cfs, err := h.cashflowSrv.ListUserCashflows(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(cfs)
}

func (h *CashflowHandler) Delete(w http.ResponseWriter, r *http.Request) {
	cfID := chi.URLParam(r, "id")
	if cfID == "" {
		http.Error(w, "missing cashflow id", http.StatusBadRequest)
		return
	}

	if err := h.cashflowSrv.DeleteCashflow(r.Context(), cfID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CashflowHandler) TriggerSync(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Trigger sync asynchronously to not block the request
	go func() {
		if err := h.syncSrv.SyncTransactions(context.Background(), userID); err != nil {
			// error logged in service
		}
	}()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"sync started"}`))
}
