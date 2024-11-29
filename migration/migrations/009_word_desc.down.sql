ALTER TABLE IF EXISTS "word"
DROP COLUMN IF EXISTS "description";

ALTER TABLE IF EXISTS "dictionary"
ALTER COLUMN "text"
TYPE text;

ALTER TABLE IF EXISTS "dictionary"
ALTER COLUMN "pronunciation"
TYPE text;

ALTER TABLE IF EXISTS "dictionary"
ALTER COLUMN "lang_code"
TYPE text;

ALTER TABLE IF EXISTS "example"
ALTER COLUMN "text"
TYPE text;

ALTER TABLE IF EXISTS "language"
ALTER COLUMN "code"
TYPE text;

ALTER TABLE IF EXISTS "language"
ALTER COLUMN "lang"
TYPE text;

ALTER TABLE IF EXISTS "user_data"
ALTER COLUMN "max_count_words"
TYPE int8;

ALTER TABLE IF EXISTS "user_subscription"
ALTER COLUMN "user_id"
TYPE varchar USING user_id::varchar;

ALTER TABLE IF EXISTS "users"
ALTER COLUMN "name"
TYPE text;

ALTER TABLE IF EXISTS "users"
ALTER COLUMN "email"
TYPE text;

ALTER TABLE IF EXISTS "users"
ALTER COLUMN "password_hash"
TYPE text;

ALTER TABLE IF EXISTS "users"
ALTER COLUMN "role"
TYPE text;

ALTER TABLE IF EXISTS "vocabulary"
ALTER COLUMN "name"
TYPE text;

ALTER TABLE IF EXISTS "vocabulary"
ALTER COLUMN "native_lang"
TYPE text;

ALTER TABLE IF EXISTS "vocabulary"
ALTER COLUMN "translate_lang"
TYPE text;