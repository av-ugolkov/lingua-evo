INSERT INTO
    users (id, name, email, password_hash, role, last_visit_at, created_at)
VALUES
    (
        '23cc06c9-73d3-40b1-9e8a-6f80db183e7a',
        'admin',
        'ugolkov.prog@gmail.com',
        '$2a$11$oP15pJXtp2ErbHWvGN05ouiMphIzrf8yXJEHkmtf.25JzgWFRaQO6',
        'admin',
        now(),
        now()
    ) ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en', 'English'),
    ('ru', 'Russian'),
    ('fi', 'Finnish'),
    ('fr', 'French'),
    ('es', 'Spanish'),
    ('it', 'Italian'),
    ('de', 'German'),
    ('pt', 'Portuguese'),
    ('sv', 'Swedish') ON CONFLICT
DO NOTHING;

CREATE TABLE IF NOT EXISTS
    "dictionary_en" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_en__text" ON "dictionary_en" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_ru" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_ru__text" ON "dictionary_ru" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_fi" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_fi__text" ON "dictionary_fi" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_fr" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_fr__text" ON "dictionary_fr" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_it" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_it__text" ON "dictionary_it" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_es" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_es__text" ON "dictionary_es" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_de" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_de__text" ON "dictionary_de" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_pt" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_pt__text" ON "dictionary_pt" ("text");

CREATE TABLE IF NOT EXISTS
    "dictionary_sv" () INHERITS ("dictionary");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_sv__text" ON "dictionary_sv" ("text");

CREATE TABLE IF NOT EXISTS
    "example_en" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_en__text" ON "example_en" ("text");

CREATE TABLE IF NOT EXISTS
    "example_ru" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_ru__text" ON "example_ru" ("text");

CREATE TABLE IF NOT EXISTS
    "example_fi" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_fi__text" ON "example_fi" ("text");

CREATE TABLE IF NOT EXISTS
    "example_fr" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_fr__text" ON "example_fr" ("text");

CREATE TABLE IF NOT EXISTS
    "example_es" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_es__text" ON "example_es" ("text");

CREATE TABLE IF NOT EXISTS
    "example_it" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_it__text" ON "example_it" ("text");

CREATE TABLE IF NOT EXISTS
    "example_de" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_de__text" ON "example_de" ("text");

CREATE TABLE IF NOT EXISTS
    "example_pt" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_pt__text" ON "example_pt" ("text");

CREATE TABLE IF NOT EXISTS
    "example_sv" () INHERITS ("example");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_sv__text" ON "example_sv" ("text");