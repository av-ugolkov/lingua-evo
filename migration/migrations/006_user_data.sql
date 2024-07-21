-- +goose Up
CREATE TABLE IF NOT EXISTS
    "subscriptions" (
        "id" BIGINT PRIMARY KEY,
        "name" VARCHAR NOT NULL,
        "description" VARCHAR,
        "add_words" BIGINT NOT NULL,
        "price_eur" BIGINT NOT NULL,
        "price_usd" BIGINT NOT NULL,
        "price_rub" BIGINT NOT NULL,
        "duration" BIGINT NOT NULL,
        "is_active" BOOLEAN DEFAULT 'false',
        "started_at" TIMESTAMP NOT NULL,
        "ended_at" TIMESTAMP NOT NULL
    );

CREATE TABLE IF NOT EXISTS
    "user_subscription" (
        "id" UUID PRIMARY KEY,
        "user_id" VARCHAR NOT NULL,
        "subscription_id" BIGINT REFERENCES "subscriptions" ("id") ON DELETE CASCADE,
        "started_at" TIMESTAMP NOT NULL,
        "ended_at" TIMESTAMP NOT NULL
    );

CREATE TABLE IF NOT EXISTS
    "user_data" (
        "user_id" UUID PRIMARY KEY REFERENCES "users" ("id") ON DELETE CASCADE,
        "max_count_words" BIGINT NOT NULL DEFAULT 300,
        "newsletter" BOOLEAN NOT NULL DEFAULT TRUE
    );

INSERT INTO
    "user_data" ("user_id")
SELECT
    "id"
FROM
    "users";

CREATE TABLE IF NOT EXISTS
    "subscribers" (
        "user_id" UUID PRIMARY KEY REFERENCES "users" ("id") ON DELETE CASCADE,
        "subscribers_id" UUID NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_user_subscribers__user_id_subscribers_id" ON "user_subscribers" ("user_id", "subscribers_id");

-- +goose Down
DROP TABLE IF EXISTS "subscriptions";

DROP TABLE IF EXISTS "user_subscription";

DROP TABLE IF EXISTS "user_subscribers";

DROP TABLE IF EXISTS "user_data";

DROP INDEX IF EXISTS "idx_unique_user_subscribers__user_id_subscribers_id";