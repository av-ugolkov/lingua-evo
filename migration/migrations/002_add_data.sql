-- +goose Up
INSERT INTO
    users (name, email, password_hash, role)
VALUES
    ('admin', 'makedonskiy07@gmail.com', '$2a$14$/55EnnJAv.3XYeKwwU6WAuAKkfO/GvSOYwADH0JSqsX2nCE2OtB.2', 'admin') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en-US', 'USA') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en-GB', 'English') ON CONFLICT
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
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get by', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get away', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get into', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get together', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get away with it', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get up to', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get up', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get rid of', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get to', '', 'en-GB', now());

INSERT INTO
    "word_en-GB" (id, text, pronunciation, lang_code, created_at)
VALUES
    (gen_random_uuid (), 'get over', '', 'en-GB', now());

-- +goose Down
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

TRUNCATE TABLE language RESTART IDENTITY CASCADE;

TRUNCATE TABLE word RESTART IDENTITY CASCADE;