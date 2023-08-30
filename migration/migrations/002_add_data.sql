-- +goose Up
INSERT INTO
    users (name, email, password_hash)
VALUES
    (
        'admin',
        'makedonskiy07@gmail.com',
        '$2a$14$/55EnnJAv.3XYeKwwU6WAuAKkfO/GvSOYwADH0JSqsX2nCE2OtB.2'
    ) ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en_US', 'USA') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en_EN', 'English') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('ru_RU', 'Russian') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('fi_FI', 'Finnish') ON CONFLICT
DO NOTHING;

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get by',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get away',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get into',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get together',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get away with it',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get up to',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get up',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get rid of',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get to',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

INSERT INTO
    word (text, lang_id)
VALUES
    (
        'get over',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_EN'
        )
    );

-- +goose Down
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

TRUNCATE TABLE language RESTART IDENTITY CASCADE;

TRUNCATE TABLE word RESTART IDENTITY CASCADE;