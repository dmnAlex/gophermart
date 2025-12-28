CREATE TABLE withdrawals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    number TEXT NOT NULL,
    sum NUMERIC(15, 2) NOT NULL CHECK (sum > 0),
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_withdrawals_user_id_processed_at ON withdrawals(user_id, processed_at DESC);