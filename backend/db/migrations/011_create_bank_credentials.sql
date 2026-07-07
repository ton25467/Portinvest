-- +goose Up
CREATE TABLE bank_credentials (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    bank_name VARCHAR(50) NOT NULL, -- e.g., 'KBANK', 'KTB', 'SCBAM'
    username VARCHAR(255) NOT NULL,
    password_encrypted VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE (user_id, bank_name)
);

-- +goose Down
DROP TABLE IF EXISTS bank_credentials;
