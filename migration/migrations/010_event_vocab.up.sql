CREATE TABLE IF NOT EXISTS
    "events_type" ("id" serial PRIMARY KEY, "name" text UNIQUE NOT NULL);

INSERT INTO
    "events_type" ("name")
VALUES
    ('vocab_created'),
    ('vocab_deleted'),
    ('vocab_updated'),
    ('vocab_renamed'),
    ('vocab_word_created'),
    ('vocab_word_deleted'),
    ('vocab_word_updated'),
    ('vocab_word_renamed') ON CONFLICT
DO NOTHING;

CREATE TABLE IF NOT EXISTS
    "events" (
        "id" uuid PRIMARY KEY,
        "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "type" bigint REFERENCES "events_type" ("id") ON DELETE CASCADE,
        "payload" jsonb NOT NULL,
        "created_at" timestamp NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_events__user_id_type_payload" ON "events" ("user_id", "type", "payload");

CREATE TABLE IF NOT EXISTS
    "events_watched" (
        "event_id" uuid REFERENCES "events" ("id") ON DELETE CASCADE,
        "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "watched_at" timestamp NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_events_watched__event_id_user_id" ON "events_watched" ("event_id", "user_id")