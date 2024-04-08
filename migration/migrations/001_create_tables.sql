-- +goose Up
CREATE TABLE IF NOT EXISTS
    "users" (
        "id" UUID PRIMARY KEY,
        "name" TEXT NOT NULL,
        "email" TEXT NOT NULL,
        "password_hash" TEXT NOT NULL,
        "role" TEXT NOT NULL,
        "last_visit_at" TIMESTAMP NOT NULL,
        "created_at" TIMESTAMP NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_users__name" ON "users" ("name");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_users__email" ON "users" ("email");

CREATE TABLE IF NOT EXISTS
    "language" ("code" TEXT PRIMARY KEY, "lang" TEXT NOT NULL);

CREATE TABLE IF NOT EXISTS
    "dictionary" (
        "id" UUID PRIMARY KEY,
        "text" TEXT NOT NULL,
        "pronunciation" TEXT,
        "lang_code" TEXT REFERENCES "language" ("code") ON DELETE CASCADE,
        "moderator" UUID,
        "created_at" TIMESTAMP NOT NULL
    );

create table
    "dictionary_en" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_en__text" ON "dictionary_en" ("text");

create table
    "dictionary_ru" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_ru__text" ON "dictionary_ru" ("text");

CREATE TABLE IF NOT EXISTS
    "example" ("id" UUID PRIMARY KEY, "text" TEXT);

CREATE TABLE IF NOT EXISTS
    "example_en" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_en__text" ON "example_en" ("text");

CREATE TABLE IF NOT EXISTS
    "example_ru" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_ru__text" ON "example_ru" ("text");

CREATE TABLE IF NOT EXISTS
    "vocabulary" (
        "id" UUID PRIMARY KEY,
        "user_id" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
        "name" TEXT NOT NULL,
        "native_lang" TEXT REFERENCES "language" ("code") ON DELETE CASCADE,
        "translate_lang" TEXT REFERENCES "language" ("code") ON DELETE CASCADE,
        "created_at" TIMESTAMP NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_vocabulary__user_id_name" ON "vocabulary" ("user_id", "name");

CREATE TABLE IF NOT EXISTS
    "tag" ("id" UUID PRIMARY KEY, "text" TEXT NOT NULL);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_tag__text" ON "tag" ("text");

CREATE TABLE IF NOT EXISTS
    "vocabulary_tag" ("vocabulary_id" UUID REFERENCES "vocabulary" ("id") ON DELETE CASCADE, "tag_id" UUID REFERENCES "tag" ("id"));

CREATE TABLE IF NOT EXISTS
    "word" (
        "id" UUID PRIMARY KEY,
        "vocabulary_id" UUID REFERENCES "vocabulary" ("id") ON DELETE CASCADE,
        "native_id" UUID REFERENCES "dictionary" ("id"),
        "created_at" TIMESTAMP NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_word__vocabulary_id_native_id" ON "word" ("vocabulary_id", "native_id");

CREATE TABLE IF NOT EXISTS
    "word_tag" ("word_id" UUID REFERENCES "word" ("id") ON DELETE CASCADE, "tag_id" UUID REFERENCES "tag" ("id"));

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_word_tag__word_id_tag_id" ON "word_tag" ("word_id", "tag_id");

CREATE TABLE IF NOT EXISTS
    "word_translate" ("word_id" UUID REFERENCES "word" ("id") ON DELETE CASCADE, "dictionary_id" UUID REFERENCES "dictionary" ("id"));

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_word_translate__word_id_dictionary_id" ON "word_translate" ("word_id", "dictionary_id");

CREATE TABLE IF NOT EXISTS
    "word_example" ("word_id" UUID REFERENCES "word" ("id") ON DELETE CASCADE, "example_id" UUID REFERENCES "example" ("id"));

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_word_example__word_id_example_id" ON "word_example" ("word_id", "example_id");

-- +goose Down
DROP TABLE IF EXISTS "users";

DROP TABLE IF EXISTS "language";

DROP TABLE IF EXISTS "dictionary";

DROP TABLE IF EXISTS "example";

DROP TABLE IF EXISTS "word";

DROP TABLE IF EXISTS "vocabulary";

DROP TABLE IF EXISTS "tag";

DROP TABLE IF EXISTS "word_tag";

DROP TABLE IF EXISTS "word_translate";

DROP TABLE IF EXISTS "word_example";