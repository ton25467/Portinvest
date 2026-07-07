package service

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"portinves/internal/domain"
	"portinves/internal/websocket"
)

type MonitoringService struct {
	serverSrv *ServerService
	hub       *websocket.Hub
	client    *http.Client
}

func NewMonitoringService(serverSrv *ServerService, hub *websocket.Hub) *MonitoringService {
	return &MonitoringService{
		serverSrv: serverSrv,
		hub:       hub,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *MonitoringService) Start(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second) // check every 15s
	slog.Info("Starting monitoring service background worker")

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping monitoring service background worker")
			return
		case <-ticker.C:
			s.runChecks(ctx)
		}
	}
}

func (s *MonitoringService) runChecks(ctx context.Context) {
	checks, err := s.serverSrv.ListActiveServiceChecks(ctx)
	if err != nil {
		slog.Error("Failed to list active service checks", "error", err)
		return
	}

	for _, check := range checks {
		go s.executeCheck(ctx, check)
	}
}

func (s *MonitoringService) executeCheck(ctx context.Context, check domain.ServiceCheck) {
	// Create check context with timeout
	checkCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(checkCtx, string(check.Method), check.Endpoint, nil)
	if err != nil {
		s.logAndBroadcast(ctx, check, domain.CheckStatusDown, 0, 0, err.Error())
		return
	}

	start := time.Now()
	resp, err := s.client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		s.logAndBroadcast(ctx, check, domain.CheckStatusDown, 0, int(elapsed.Milliseconds()), err.Error())
		return
	}
	defer resp.Body.Close()

	status := domain.CheckStatusUp
	if resp.StatusCode != check.ExpectedStatus {
		status = domain.CheckStatusDown
	}

	s.logAndBroadcast(ctx, check, status, resp.StatusCode, int(elapsed.Milliseconds()), "")
}

type ServiceCheckWSMessage struct {
	Type           string              `json:"type"` // e.g. "service_check_result"
	ServerID       string              `json:"server_id"`
	ServiceCheckID string              `json:"service_check_id"`
	Status         domain.CheckStatus  `json:"status"`
	StatusCode     int                 `json:"status_code"`
	ResponseTimeMs int                 `json:"response_time_ms"`
	CheckedAt      time.Time           `json:"checked_at"`
	ServerStatus   domain.ServerStatus `json:"server_status"`
}

func (s *MonitoringService) logAndBroadcast(ctx context.Context, check domain.ServiceCheck, status domain.CheckStatus, statusCode, responseTime int, errMsg string) {
	log, err := s.serverSrv.LogCheckResult(ctx, check.ID, status, statusCode, responseTime, errMsg)
	if err != nil {
		slog.Error("Failed to log check result", "check_id", check.ID, "error", err)
		return
	}

	// Fetch server to get its updated status
	srv, err := s.serverSrv.GetServer(ctx, check.ServerID)
	srvStatus := domain.ServerStatusOnline
	if err == nil {
		srvStatus = srv.Status
	}

	msg := ServiceCheckWSMessage{
		Type:           "service_check_result",
		ServerID:       check.ServerID,
		ServiceCheckID: check.ID,
		Status:         log.Status,
		StatusCode:     log.StatusCode,
		ResponseTimeMs: log.ResponseTimeMs,
		CheckedAt:      log.CheckedAt,
		ServerStatus:   srvStatus,
	}

	s.hub.Broadcast(msg)
}
