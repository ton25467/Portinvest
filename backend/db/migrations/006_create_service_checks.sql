-- +goose Up
CREATE TABLE service_checks (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    server_id        UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
    name             VARCHAR(255) NOT NULL,
    endpoint         TEXT NOT NULL,
    method           VARCHAR(10) NOT NULL DEFAULT 'GET' CHECK (method IN ('GET','POST')),
    expected_status  INTEGER NOT NULL DEFAULT 200,
    interval_seconds INTEGER NOT NULL DEFAULT 60,
    is_active        BOOLEAN NOT NULL DEFAULT true,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_service_checks_server_id ON service_checks(server_id);
CREATE INDEX idx_service_checks_active ON service_checks(is_active);

-- +goose Down
DROP TABLE IF EXISTS service_checks;
