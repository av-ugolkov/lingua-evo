-- +goose Up
CREATE TABLE IF NOT EXISTS
    "session" (
        id UUID PRIMARY KEY,
        user_id UUID REFERENCES users (id) ON DELETE CASCADE,
        refresh_token UUID NOT NULL,
        expires_in BIGINT NOT NULL,
        "created_at" timestamp with time zone NOT NULL
    );

-- +goose Down
DROP TABLE IF EXISTS "session";