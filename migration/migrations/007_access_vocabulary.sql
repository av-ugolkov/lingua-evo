-- +goose Up
CREATE TABLE IF NOT EXISTS
    "access" ("id" BIGINT PRIMARY KEY, "type" VARCHAR(255) NOT NULL UNIQUE);

INSERT INTO
    "access" ("id", "type")
VALUES
    (0, 'private'),
    (1, 'some_one'),
    (2, 'subscribers'),
    (3, 'public');

ALTER TABLE IF EXISTS "vocabulary"
ADD COLUMN IF NOT EXISTS "access" BIGINT NOT NULL DEFAULT 2;

ALTER TABLE IF EXISTS "vocabulary"
ADD COLUMN IF NOT EXISTS "access_edit" BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS
    "vocabulary_users_access" (
        "vocab_id" UUID REFERENCES "vocabulary" ("id") ON DELETE CASCADE,
        "subscriber_id" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
        "editor" BOOLEAN NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_vocabulary_users_access__vocab_id_subscriber_id" ON "vocabulary_users_access" ("vocab_id", "subscriber_id");

-- +goose Down
DROP TABLE IF EXISTS "access";

ALTER TABLE IF EXISTS "vocabulary"
DROP COLUMN IF EXISTS "access";

DROP TABLE IF EXISTS "vocabulary_users_access";

DROP INDEX IF EXISTS "idx_unique_vocabulary_users_access__vocab_id_subscriber_id";