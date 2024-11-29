ALTER TABLE IF EXISTS "users"
RENAME COLUMN "name" TO "nickname";

ALTER TABLE IF EXISTS "users"
RENAME COLUMN "last_visit_at" TO "visited_at";

ALTER TABLE IF EXISTS "user_data"
ADD COLUMN IF NOT EXISTS "name" varchar(100) default '';

ALTER TABLE IF EXISTS "user_data"
ADD COLUMN IF NOT EXISTS "surname" varchar(100) default '';

ALTER TABLE IF EXISTS "users"
ADD COLUMN IF NOT EXISTS "max_count_words" BIGINT NOT null default 300;

ALTER TABLE IF EXISTS "user_data"
DROP COLUMN IF EXISTS "max_count_words";

ALTER TABLE IF EXISTS "user_data"
DROP COLUMN IF EXISTS "newsletter";

CREATE TABLE IF NOT EXISTS
    "user_newsletters" ("user_id" UUID REFERENCES "users" ("id") ON DELETE CASCADE, "news" BOOLEAN NOT NULL default true);

INSERT INTO
    "user_newsletters" ("user_id", "news")
SELECT
    "id",
    true
FROM
    "users" ON CONFLICT
DO NOTHING;