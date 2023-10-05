-- +goose Up
INSERT INTO
    users (name, email, password_hash)
VALUES
    ('admin', 'makedonskiy07@gmail.com', '$2a$14$/55EnnJAv.3XYeKwwU6WAuAKkfO/GvSOYwADH0JSqsX2nCE2OtB.2') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en_us', 'USA') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('en_gb', 'English') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('ru_ru', 'Russian') ON CONFLICT
DO NOTHING;

INSERT INTO
    language (code, lang)
VALUES
    ('fi_fi', 'Finnish') ON CONFLICT
DO NOTHING;

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get by',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get away',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get into',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get together',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get away with it',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get up to',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get up',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get rid of',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get to',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

INSERT INTO
    word_en_gb (id, text, lang_id)
VALUES
    (
        gen_random_uuid (),
        'get over',
        (
            SELECT
                id
            FROM
                language
            where
                code = 'en_gb'
        )
    );

-- +goose Down
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

TRUNCATE TABLE language RESTART IDENTITY CASCADE;

TRUNCATE TABLE word RESTART IDENTITY CASCADE;