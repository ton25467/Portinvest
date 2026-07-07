-- +goose Up
CREATE TABLE servers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(255) NOT NULL,
    host            VARCHAR(255) NOT NULL,
    port            INTEGER NOT NULL,
    type            VARCHAR(20) NOT NULL CHECK (type IN ('web','db','api','cache')),
    status          VARCHAR(20) NOT NULL DEFAULT 'offline' CHECK (status IN ('online','offline','degraded')),
    last_checked_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_servers_user_id ON servers(user_id);
CREATE INDEX idx_servers_status ON servers(status);

-- +goose Down
DROP TABLE IF EXISTS servers;
