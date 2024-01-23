-- +goose Up
CREATE TABLE IF NOT EXISTS
    users (
        id UUID PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        password_hash TEXT NOT NULL,
        role TEXT NOT NULL,
        last_visit_at TIMESTAMP NOT NULL,
        created_at TIMESTAMP NOT NULL
    );

CREATE UNIQUE INDEX IF not EXISTS idx_unique_users__name ON users (name);

CREATE UNIQUE INDEX IF not EXISTS idx_unique_users__email ON users (email);

CREATE TABLE IF NOT EXISTS
    language (code TEXT NOT NULL PRIMARY KEY, lang TEXT NOT NULL);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_language__code ON language (code);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_language__code_lang ON language (code, lang);

CREATE TABLE IF NOT EXISTS
    word (
        id UUID PRIMARY KEY,
        text TEXT NOT NULL,
        pronunciation TEXT,
        lang_code TEXT,
        created_at TIMESTAMP NOT NULL,
        CONSTRAINT word_lang_code_fkey FOREIGN KEY (lang_code) REFERENCES language (code) ON DELETE CASCADE
    );

create table
    "word_en" () INHERITS (word);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_word_en__text" ON "word_en" (text);

create table
    word_ru () INHERITS (word);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_word_ru__text" ON "word_ru" (text);

CREATE TABLE IF NOT EXISTS
    example (id UUID DEFAULT gen_random_uuid () PRIMARY KEY, text TEXT);

CREATE TABLE IF NOT EXISTS
    "example_en" () INHERITS (example);

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_en__text" ON "example_en" (text);

CREATE TABLE IF NOT EXISTS
    example_ru () INHERITS (example);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_example_ru__text ON example_ru (text);

CREATE TABLE IF NOT EXISTS
    dictionary (id UUID PRIMARY KEY, user_id UUID REFERENCES users (id) NOT NULL, name TEXT NOT NULL, tags UUID[] NOT NULL);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_dictionary__user_id_name ON dictionary (user_id, name);

CREATE TABLE IF NOT EXISTS
    tag (id UUID PRIMARY KEY, text TEXT NOT NULL);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_tag__text ON tag (text);

CREATE TABLE IF NOT EXISTS
    vocabulary (
        dictionary_id UUID REFERENCES dictionary (id) NOT NULL,
        native_word UUID NOT NULL,
        translate_word UUID[] NOT NULL,
        examples UUID[],
        tags INT[]
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_vocabulary__dictionary_id_native_word ON vocabulary (dictionary_id, native_word);

-- +goose Down
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS language;

DROP TABLE IF EXISTS example;

DROP TABLE IF EXISTS word;

DROP TABLE IF EXISTS dictionary;

DROP TABLE IF EXISTS vocabulary;

DROP TABLE IF EXISTS tag;