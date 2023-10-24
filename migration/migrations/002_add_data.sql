-- +goose Up
INSERT INTO
    users (id, name, email, password_hash, role, last_visit_at, created_at)
VALUES
    (
        gen_random_uuid (),
        'admin',
        'makedonskiy07@gmail.com',
        '$2a$14$/55EnnJAv.3XYeKwwU6WAuAKkfO/GvSOYwADH0JSqsX2nCE2OtB.2',
        'admin',
        now(),
        now()
    ) ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en-US', 'USA') ON CONFLICT
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

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get by', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get away', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get into', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get together', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get away with it', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get up to', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get up', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get rid of', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get to', '', 'en', now());

INSERT INTO
    "word_en" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get over', '', 'en', now());

-- +goose Down
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

TRUNCATE TABLE language RESTART IDENTITY CASCADE;

TRUNCATE TABLE word RESTART IDENTITY CASCADE;