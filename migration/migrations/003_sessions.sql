-- +goose Up
CREATE TABLE IF NOT EXISTS
    "session" (
        id UUID PRIMARY KEY,
        user_id UUID REFERENCES users (id) ON DELETE CASCADE,
        refresh_token UUID NOT NULL,
        expires_at TIMESTAMP NOT NULL,
        "created_at" TIMESTAMP NOT NULL
    );

-- +goose Down
DROP TABLE IF EXISTS "session";