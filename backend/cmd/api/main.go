package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"portinves/config"
	"portinves/internal/handler"
	"portinves/internal/repository/postgres"
	"portinves/internal/service"
	"portinves/internal/websocket"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Database Migrations
	sqlDB, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		slog.Error("Failed to open SQL connection for migrations", "error", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	slog.Info("Running database migrations...")
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to set goose dialect", "error", err)
		os.Exit(1)
	}

	migrationDir := "db/migrations"
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		migrationDir = "../../db/migrations"
		if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
			migrationDir = "../db/migrations"
		}
	}

	if err := goose.Up(sqlDB, migrationDir); err != nil {
		slog.Error("Database migration failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migrations completed successfully")

	// 2. pgx Connection Pool
	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL, cfg.MaxDBConns)
	if err != nil {
		slog.Error("Failed to create pgx connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 3. Repositories
	userRepo := postgres.NewUserRepo(pool)
	portfolioRepo := postgres.NewPortfolioRepo(pool)
	holdingRepo := postgres.NewHoldingRepo(pool)
	txRepo := postgres.NewTransactionRepo(pool)
	serverRepo := postgres.NewServerRepo(pool)
	checkRepo := postgres.NewServiceCheckRepo(pool)
	logRepo := postgres.NewUptimeLogRepo(pool)
	dashboardRepo := postgres.NewDashboardRepo(pool)
	cashflowRepo := postgres.NewCashflowRepo(pool)
	notifRepo := postgres.NewNotificationRepo(pool)
	credRepo := postgres.NewBankCredentialRepo(pool)

	// 4. Services & Websocket Hub
	authSrv := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	portfolioSrv := service.NewPortfolioService(pool, portfolioRepo, holdingRepo, txRepo)
	serverSrv := service.NewServerService(serverRepo, checkRepo, logRepo)

	hub := websocket.NewHub()

	notifSrv := service.NewNotificationService(notifRepo, hub)
	cashflowSrv := service.NewCashflowService(cashflowRepo, notifSrv)
	syncSrv := service.NewWebScraperSyncService(cashflowSrv, notifSrv, credRepo)
	defer syncSrv.Stop()

	monitoringSrv := service.NewMonitoringService(serverSrv, hub)
	go monitoringSrv.Start(ctx)

	priceSrv := service.NewPriceService(portfolioSrv)
	go priceSrv.Start(ctx)

	// 5. Handlers
	authH := handler.NewAuthHandler(authSrv)
	portfolioH := handler.NewPortfolioHandler(portfolioSrv)
	holdingH := handler.NewHoldingHandler(portfolioSrv)
	serverH := handler.NewServerHandler(serverSrv)
	dashboardH := handler.NewDashboardHandler(dashboardRepo)
	wsH := handler.NewWSHandler(hub, authSrv)
	notifH := handler.NewNotificationHandler(notifSrv)
	cashflowH := handler.NewCashflowHandler(cashflowSrv, syncSrv)
	credH := handler.NewBankCredentialHandler(credRepo)

	// 6. Router Setup
	r := handler.NewRouter(
		authSrv,
		authH,
		portfolioH,
		holdingH,
		serverH,
		dashboardH,
		wsH,
		notifH,
		cashflowH,
		credH,
		cfg.CORSOrigins,
	)

	// Start server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Server is running", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Server is shutting down...")
	cancel() // stops monitoring background workers

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
	}
	slog.Info("Server exited gracefully")
}
