ALTER TABLE IF EXISTS "users"
DROP COLUMN IF EXISTS "max_count_words";

ALTER TABLE IF EXISTS "users"
ALTER COLUMN "password_hash"
DROP NOT NULL;

ALTER TABLE IF EXISTS "users"
ADD COLUMN IF NOT EXISTS "google_id" TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_users__google_id" ON "users" ("google_id");

ALTER TABLE IF EXISTS "user_data"
ADD COLUMN IF NOT EXISTS "url_avatar" TEXT NOT NULL DEFAULT '';