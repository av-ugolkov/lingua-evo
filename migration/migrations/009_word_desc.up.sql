ALTER TABLE word
ADD COLUMN IF NOT EXISTS description varchar(255) default '';

ALTER TABLE dictionary
ALTER COLUMN text
TYPE varchar(100);

ALTER TABLE dictionary
ALTER COLUMN pronunciation
TYPE varchar(100);

ALTER TABLE dictionary
ALTER COLUMN lang_code
TYPE varchar(10);

ALTER TABLE example
ALTER COLUMN text
TYPE varchar(255);

ALTER TABLE language
ALTER COLUMN code
TYPE varchar(10);

ALTER TABLE language
ALTER COLUMN lang
TYPE varchar(50);

ALTER TABLE user_data
ALTER COLUMN max_count_words
TYPE smallint;

ALTER TABLE user_subscription
ALTER COLUMN user_id
TYPE uuid USING user_id::uuid;

ALTER TABLE users
ALTER COLUMN name
TYPE varchar(50);

ALTER TABLE users
ALTER COLUMN email
TYPE varchar(50);

ALTER TABLE users
ALTER COLUMN password_hash
TYPE varchar(255);

ALTER TABLE users
ALTER COLUMN role
TYPE varchar(25);

ALTER TABLE vocabulary
ALTER COLUMN name
TYPE varchar(50);

ALTER TABLE vocabulary
ALTER COLUMN native_lang
TYPE varchar(25);

ALTER TABLE vocabulary
ALTER COLUMN translate_lang
TYPE varchar(25);