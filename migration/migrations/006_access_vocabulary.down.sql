ALTER TABLE IF EXISTS "vocabulary"
DROP COLUMN IF EXISTS "access";

DROP TABLE IF EXISTS "vocabulary_users_access" CASCADE;

DROP INDEX IF EXISTS "idx_unique_vocabulary_users_access__vocab_id_subscriber_id";

DROP TABLE IF EXISTS "access" CASCADE;