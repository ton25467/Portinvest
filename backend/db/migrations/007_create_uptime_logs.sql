-- +goose Up
CREATE TABLE uptime_logs (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    service_check_id  UUID NOT NULL REFERENCES service_checks(id) ON DELETE CASCADE,
    status            VARCHAR(20) NOT NULL CHECK (status IN ('up','down','degraded')),
    status_code       INTEGER NOT NULL DEFAULT 0,
    response_time_ms  INTEGER NOT NULL DEFAULT 0,
    error_message     TEXT NOT NULL DEFAULT '',
    checked_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_uptime_logs_service_check_id ON uptime_logs(service_check_id);
CREATE INDEX idx_uptime_logs_checked_at ON uptime_logs(checked_at);

-- +goose Down
DROP TABLE IF EXISTS uptime_logs;
