-- +goose Up
CREATE TABLE IF NOT EXISTS
    users (
        id UUID DEFAULT gen_random_uuid () PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT NOT NULL,
        password_hash TEXT NOT NULL,
        last_visit date
    );

CREATE UNIQUE INDEX IF not EXISTS idx_unique_users__name ON users (name);

CREATE UNIQUE INDEX IF not EXISTS idx_unique_users__email ON users (email);

CREATE TABLE IF NOT EXISTS
    language (id bigserial PRIMARY KEY, code TEXT NOT NULL, lang TEXT NOT NULL);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_language__lang_code ON language (lang, code);

CREATE TABLE IF NOT EXISTS
    word (
        id UUID NOT NULL PRIMARY KEY,
        text TEXT NOT NULL,
        pronunciation TEXT,
        lang_id INT,
        created_at TIMESTAMP NOT NULL,
        CONSTRAINT word_lang_id_fkey FOREIGN KEY (lang_id) REFERENCES language (id) ON DELETE CASCADE
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_word__text_lang_id ON word (text, lang_id);

create table
    "word_en-GB" () inherits (word);

create table
    word_ru () inherits (word);

CREATE TABLE IF NOT EXISTS
    example (id UUID DEFAULT gen_random_uuid () PRIMARY KEY, example TEXT);

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_example__example ON example (example);

CREATE TABLE IF NOT EXISTS
    dictionary (
        user_id UUID REFERENCES users (id) NOT NULL,
        original_word UUID REFERENCES word (id) NOT NULL,
        translate_word UUID[] NOT NULL,
        examples UUID[],
        tags INT[]
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_dictionary__user_id_original_word ON dictionary (user_id, original_word);

CREATE TABLE IF NOT EXISTS
    tag (id bigserial PRIMARY KEY, tag TEXT NOT NULL);

-- +goose Down
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS language;

DROP TABLE IF EXISTS example;

DROP TABLE IF EXISTS word;

DROP TABLE IF EXISTS dictionary;

DROP TABLE IF EXISTS tag;