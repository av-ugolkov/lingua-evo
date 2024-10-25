CREATE TABLE IF NOT EXISTS
    "event" ("id" uuid PRIMARY KEY, "payload" jsonb NOT NULL, "created_at" timestamp NOT NULL);

CREATE INDEX IF NOT EXISTS "idx_hash_event__id" ON "event" USING HASH ("id");

CREATE TABLE IF NOT EXISTS
    "event_vocab" ("vocab_id" uuid NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE) INHERITS ("event");

CREATE INDEX IF NOT EXISTS "idx_hash_dictionary__word_id" ON "dictionary" USING HASH ("id");

CREATE TABLE IF NOT EXISTS
    "event_watched" (
        "event_id" uuid REFERENCES "event" ("id") ON DELETE CASCADE,
        "user_id" uuid REFERENCES "users" ("id") ON DELETE CASCADE,
        "watched_at" timestamp NOT NULL
    );