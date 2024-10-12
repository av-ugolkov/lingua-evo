ALTER TABLE IF EXISTS "word"
DROP COLUMN IF EXISTS "pronunciation";

DROP TABLE IF EXISTS "user_subscription" CASCADE;

DROP TABLE IF EXISTS "user_subscribers" CASCADE;

DROP TABLE IF EXISTS "subscriptions" CASCADE;

DROP TABLE IF EXISTS "user_data" CASCADE;

DROP TABLE IF EXISTS "subscribers" CASCADE;

DROP INDEX IF EXISTS "idx_unique_user_subscribers__user_id_subscribers_id";