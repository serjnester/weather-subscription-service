-- +goose Up
CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL,
    city TEXT NOT NULL,
    frequency TEXT NOT NULL CHECK (frequency IN ('hourly', 'daily')),
    confirmed BOOLEAN NOT NULL DEFAULT FALSE,
    token TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(email, city)
);

-- +goose Down
DROP TABLE IF EXISTS subscriptions;
