-- +goose Up
-- Seed default user (password is 'password')
INSERT INTO users (id, email, password_hash, name, role)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'user@example.com',
    '$2a$10$3ElmzhAH05vEHtjx5ujcuettLgzzT/dq52ly3Lwx56sIwSi0qH26K',
    'Demo User',
    'user'
) ON CONFLICT (email) DO NOTHING;

-- Seed default THB portfolio
INSERT INTO portfolios (id, user_id, name, description, currency)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    'My Investment Portfolio',
    'Stock and Mutual Fund Holdings',
    'THB'
) ON CONFLICT (id) DO NOTHING;

-- Seed user's actual stock and fund holdings
INSERT INTO holdings (id, portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price)
VALUES
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'TISCO', 'Tisco Financial Group', 'stock', 100, 98.50, 99.00),
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'KTB', 'Krung Thai Bank', 'stock', 500, 18.20, 18.60),
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'K-JPX-A(A)', 'K Japan Share Index Fund', 'etf', 1000, 12.40, 12.80),
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'K-US500X-A(A)', 'K US Equity Index Fund (S&P 500)', 'etf', 2000, 15.60, 16.10),
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'K-USXNDQ-A(A)', 'K US Nasdaq 100 Index Fund', 'etf', 1500, 22.30, 23.40),
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'K-GOLD-A(A)', 'K Gold Fund', 'etf', 800, 10.50, 11.20),
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'K-WPBALANCED', 'K Wealth Plus Balanced Fund', 'etf', 3000, 10.00, 10.15),
    (uuid_generate_v4(), '00000000-0000-0000-0000-000000000002', 'SCBS&P500E', 'SCB S&P 500 Index Fund (e-class)', 'etf', 2500, 14.80, 15.30)
ON CONFLICT DO NOTHING;

-- Seed transaction history for these holdings so P&L looks realistic
-- We can fetch the holdings we just created, but since they have dynamic UUIDs, we just insert holdings directly
-- This is fine for display.

-- +goose Down
DELETE FROM holdings WHERE portfolio_id = '00000000-0000-0000-0000-000000000002';
DELETE FROM portfolios WHERE id = '00000000-0000-0000-0000-000000000002';
DELETE FROM users WHERE id = '00000000-0000-0000-0000-000000000001';
