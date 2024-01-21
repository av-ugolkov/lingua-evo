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

-- +goose Down
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

TRUNCATE TABLE language RESTART IDENTITY CASCADE;

TRUNCATE TABLE word RESTART IDENTITY CASCADE;