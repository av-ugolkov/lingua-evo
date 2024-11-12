ALTER TABLE IF EXISTS "users"
RENAME COLUMN "nickname" TO "name";

ALTER TABLE IF EXISTS "users"
RENAME COLUMN "visited_at" TO "last_visit_at";

ALTER TABLE IF EXISTS "user_data"
DROP COLUMN IF EXISTS 'name';

ALTER TABLE IF EXISTS "user_data"
DROP COLUMN IF EXISTS 'surname';

ALTER TABLE IF EXISTS "users"
DROP COLUMN IF EXISTS "max_count_words";

ALTER TABLE IF EXISTS "user_data"
ADD COLUMN IF NOT EXISTS "max_count_words" BIGINT NOT null default 300;

ALTER TABLE IF EXISTS "user_data"
ADD COLUMN IF NOT EXISTS "newsletter" BOOLEAN NOT NULL DEFAULT true;

DROP TABLE IF EXISTS "user_newsletters";