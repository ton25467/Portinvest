package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"portinves/internal/handler/middleware"
	"portinves/internal/service"
)

func NewRouter(
	authSrv *service.AuthService,
	authH *AuthHandler,
	portfolioH *PortfolioHandler,
	holdingH *HoldingHandler,
	serverH *ServerHandler,
	dashboardH *DashboardHandler,
	wsH *WSHandler,
	notifH *NotificationHandler,
	cashflowH *CashflowHandler,
	credH *BankCredentialHandler,
	corsOrigins []string,
) http.Handler {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.CORS(corsOrigins))
	r.Use(middleware.Logger)

	// Health Check
	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	})

	// Public Auth Routes
	r.Post("/api/v1/auth/register", authH.Register)
	r.Post("/api/v1/auth/login", authH.Login)

	// WebSocket Endpoint
	r.Get("/ws", wsH.ServeWS)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(authSrv))

		// Portfolios
		r.Post("/api/v1/portfolios", portfolioH.Create)
		r.Get("/api/v1/portfolios", portfolioH.List)
		r.Get("/api/v1/portfolios/{id}", portfolioH.GetByID)
		r.Get("/api/v1/portfolios/{id}/summary", portfolioH.GetSummary)

		// Holdings & Transactions
		r.Get("/api/v1/portfolios/{id}/holdings", holdingH.ListHoldings)
		r.Post("/api/v1/portfolios/{id}/transactions", holdingH.AddTransaction)

		// Servers
		r.Post("/api/v1/servers", serverH.Create)
		r.Get("/api/v1/servers", serverH.List)
		r.Get("/api/v1/servers/{id}", serverH.GetByID)
		r.Put("/api/v1/servers/{id}", serverH.Update)
		r.Delete("/api/v1/servers/{id}", serverH.Delete)
		r.Post("/api/v1/servers/{id}/checks", serverH.CreateCheck)
		r.Get("/api/v1/servers/{id}/checks", serverH.ListChecks)
		r.Get("/api/v1/servers/{id}/logs", serverH.ListLogs)

		// Dashboard Overview
		r.Get("/api/v1/dashboard/overview", dashboardH.GetOverview)

		// Cashflows
		r.Post("/api/v1/cashflows", cashflowH.Create)
		r.Get("/api/v1/cashflows", cashflowH.List)
		r.Delete("/api/v1/cashflows/{id}", cashflowH.Delete)
		r.Post("/api/v1/cashflows/sync", cashflowH.TriggerSync)

		// Notifications
		r.Get("/api/v1/notifications", notifH.List)
		r.Put("/api/v1/notifications/{id}/read", notifH.MarkAsRead)

		// Bank Credentials
		r.Post("/api/v1/settings/banks", credH.Save)
	})

	return r
}
