CREATE TABLE IF NOT EXISTS
    "access" ("id" BIGINT PRIMARY KEY, "type" VARCHAR(25) NOT NULL UNIQUE, "name" VARCHAR(25) NOT NULL UNIQUE);

INSERT INTO
    "access" ("id", "type", "name")
VALUES
    (0, 'private', 'Private'),
    (1, 'subscribers', 'Subscribers'),
    (2, 'public', 'Public') ON CONFLICT
DO NOTHING;

ALTER TABLE IF EXISTS "vocabulary"
ADD COLUMN IF NOT EXISTS "access" BIGINT NOT NULL DEFAULT 2 CONSTRAINT "vocabulary_access_fkey" REFERENCES "access" ("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS
    "vocabulary_users_access" (
        "vocab_id" UUID REFERENCES "vocabulary" ("id") ON DELETE CASCADE,
        "subscriber_id" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
        "editor" BOOLEAN NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_vocabulary_users_access__vocab_id_subscriber_id" ON "vocabulary_users_access" ("vocab_id", "subscriber_id");

ALTER TABLE IF EXISTS "vocabulary"
ADD COLUMN IF NOT EXISTS "description" VARCHAR(255) NOT NULL DEFAULT '';