ALTER TABLE IF EXISTS "users"
ADD COLUMN IF NOT EXISTS "max_count_words" BIGINT NOT null default 300;

ALTER TABLE IF EXISTS "users"
ALTER COLUMN "password_hash"
TYPE VARCHAR(255) NOT NULL;

ALTER TABLE IF EXISTS "users"
DROP COLUMN IF EXISTS "google_id" TEXT;

DROP INDEX IF EXISTS "idx_unique_users__google_id" ON "users" ("google_id");

ALTER TABLE IF EXISTS "user_data"
DROP COLUMN IF EXISTS "url_avatar" TEXT NOT NULL DEFAULT '';