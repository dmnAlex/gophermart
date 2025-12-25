CREATE TYPE order_status AS ENUM (
    'NEW',
    'PROCESSING',
    'INVALID',
    'PROCESSED'
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    number TEXT NOT NULL UNIQUE,
    status order_status NOT NULL DEFAULT 'NEW',
    accrual NUMERIC(15, 2) NULL CHECK (accrual >= 0),
    uploaded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_locked BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_orders_user_id_uploaded_at ON orders(user_id, uploaded_at DESC);