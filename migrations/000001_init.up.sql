CREATE TABLE IF NOT EXISTS subscriptions
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date   DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (end_date IS NULL OR end_date >= start_date)
);