ALTER TABLE IF EXISTS word
DROP COLUMN IF EXISTS description;

ALTER TABLE dictionary
ALTER COLUMN text
TYPE text;

ALTER TABLE dictionary
ALTER COLUMN pronunciation
TYPE text;

ALTER TABLE dictionary
ALTER COLUMN lang_code
TYPE text;

ALTER TABLE example
ALTER COLUMN text
TYPE text;

ALTER TABLE language
ALTER COLUMN code
TYPE text;

ALTER TABLE language
ALTER COLUMN lang
TYPE text;

ALTER TABLE user_data
ALTER COLUMN max_count_words
TYPE int8;

ALTER TABLE user_subscription
ALTER COLUMN user_id
TYPE varchar USING user_id::varchar;

ALTER TABLE users
ALTER COLUMN name
TYPE text;

ALTER TABLE users
ALTER COLUMN email
TYPE text;

ALTER TABLE users
ALTER COLUMN password_hash
TYPE text;

ALTER TABLE users
ALTER COLUMN role
TYPE text;

ALTER TABLE vocabulary
ALTER COLUMN name
TYPE text;

ALTER TABLE vocabulary
ALTER COLUMN native_lang
TYPE text;

ALTER TABLE vocabulary
ALTER COLUMN translate_lang
TYPE text;