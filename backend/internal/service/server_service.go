package service

import (
	"context"
	"time"

	"portinves/internal/domain"
)

type ServerService struct {
	serverRepo domain.ServerRepository
	checkRepo  domain.ServiceCheckRepository
	logRepo    domain.UptimeLogRepository
}

func NewServerService(
	serverRepo domain.ServerRepository,
	checkRepo domain.ServiceCheckRepository,
	logRepo domain.UptimeLogRepository,
) *ServerService {
	return &ServerService{
		serverRepo: serverRepo,
		checkRepo:  checkRepo,
		logRepo:    logRepo,
	}
}

func (s *ServerService) CreateServer(ctx context.Context, userID, name, host string, port int, srvType domain.ServerType) (*domain.Server, error) {
	if name == "" || host == "" {
		return nil, domain.ErrBadRequest("name and host are required")
	}
	srv := &domain.Server{
		UserID: userID,
		Name:   name,
		Host:   host,
		Port:   port,
		Type:   srvType,
		Status: domain.ServerStatusOnline, // default status
	}
	if err := s.serverRepo.Create(ctx, srv); err != nil {
		return nil, err
	}
	return srv, nil
}

func (s *ServerService) GetServer(ctx context.Context, id string) (*domain.Server, error) {
	return s.serverRepo.GetByID(ctx, id)
}

func (s *ServerService) ListServers(ctx context.Context, userID string) ([]domain.Server, error) {
	return s.serverRepo.ListByUserID(ctx, userID)
}

func (s *ServerService) UpdateServer(ctx context.Context, srv *domain.Server) error {
	return s.serverRepo.Update(ctx, srv)
}

func (s *ServerService) DeleteServer(ctx context.Context, id string) error {
	return s.serverRepo.Delete(ctx, id)
}

func (s *ServerService) CreateServiceCheck(ctx context.Context, serverID, name, endpoint string, method domain.HTTPMethod, expectedStatus, interval int) (*domain.ServiceCheck, error) {
	if name == "" || endpoint == "" {
		return nil, domain.ErrBadRequest("name and endpoint are required")
	}
	if interval <= 0 {
		interval = 60
	}
	if method == "" {
		method = domain.HTTPMethodGET
	}
	if expectedStatus == 0 {
		expectedStatus = 200
	}

	sc := &domain.ServiceCheck{
		ServerID:        serverID,
		Name:            name,
		Endpoint:        endpoint,
		Method:          method,
		ExpectedStatus:  expectedStatus,
		IntervalSeconds: interval,
		IsActive:        true,
	}
	if err := s.checkRepo.Create(ctx, sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *ServerService) ListServiceChecks(ctx context.Context, serverID string) ([]domain.ServiceCheck, error) {
	return s.checkRepo.ListByServerID(ctx, serverID)
}

func (s *ServerService) ListActiveServiceChecks(ctx context.Context) ([]domain.ServiceCheck, error) {
	return s.checkRepo.ListActive(ctx)
}

func (s *ServerService) LogCheckResult(ctx context.Context, checkID string, status domain.CheckStatus, statusCode, responseTime int, errMsg string) (*domain.UptimeLog, error) {
	log := &domain.UptimeLog{
		ServiceCheckID: checkID,
		Status:         status,
		StatusCode:     statusCode,
		ResponseTimeMs: responseTime,
		ErrorMessage:   errMsg,
		CheckedAt:      time.Now(),
	}

	if err := s.logRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	// Update the parent Server status
	check, err := s.checkRepo.GetByID(ctx, checkID)
	if err == nil {
		srv, err := s.serverRepo.GetByID(ctx, check.ServerID)
		if err == nil {
			now := time.Now()
			srv.LastCheckedAt = &now
			// Map check status to server status
			if status == domain.CheckStatusDown {
				srv.Status = domain.ServerStatusOffline
			} else if status == domain.CheckStatusDegraded {
				srv.Status = domain.ServerStatusDegraded
			} else {
				srv.Status = domain.ServerStatusOnline
			}
			_ = s.serverRepo.Update(ctx, srv)
		}
	}

	return log, nil
}

func (s *ServerService) GetUptimeLogs(ctx context.Context, serverID string, limit int) ([]domain.UptimeLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.logRepo.GetUptimeByServerID(ctx, serverID, limit)
}
