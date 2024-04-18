-- +goose Up
INSERT INTO
    users (id, name, email, password_hash, role, last_visit_at, created_at)
VALUES
    (
        gen_random_uuid (),
        'makedonskiy',
        'makedonskiy07@gmail.com',
        '$2a$11$Br7apwhrxr1yvfCSzsz3rec0m9MXRUiyBeY6V543cBXVwl/TKd.fO',
        'admin',
        now(),
        now()
    ) ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en', 'English') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('ru', 'Russian') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('fi', 'Finnish') ON CONFLICT
DO NOTHING;

CREATE TABLE
    "dictionary_en" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_en__text" ON "dictionary_en" ("text");

CREATE TABLE
    "dictionary_ru" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_ru__text" ON "dictionary_ru" ("text");

CREATE TABLE
    "dictionary_fi" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_fi__text" ON "dictionary_fi" ("text");

CREATE TABLE IF NOT EXISTS
    "example_en" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_en__text" ON "example_en" ("text");

CREATE TABLE IF NOT EXISTS
    "example_ru" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_ru__text" ON "example_ru" ("text");

CREATE TABLE IF NOT EXISTS
    "example_fi" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_fi__text" ON "example_fi" ("text");

-- +goose Down
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

TRUNCATE TABLE language RESTART IDENTITY CASCADE;

TRUNCATE TABLE word RESTART IDENTITY CASCADE;