CREATE TABLE IF NOT EXISTS
    "events" (
        "id" uuid PRIMARY KEY,
        "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "payload" jsonb NOT NULL,
        "created_at" timestamp NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_events__user_id_payload" ON "events" ("user_id", "payload");

CREATE INDEX IF NOT EXISTS "idx_hash_event__id" ON "events" USING HASH ("id");

CREATE INDEX IF NOT EXISTS "idx_hash_dictionary__word_id" ON "dictionary" USING HASH ("id");

CREATE TABLE IF NOT EXISTS
    "event_watched" (
        "event_id" uuid REFERENCES "events" ("id") ON DELETE CASCADE,
        "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "watched_at" timestamp NOT NULL
    );