CREATE TABLE IF NOT EXISTS
    "goose_db_version" (
        "id" serial4 NOT NULL,
        "version_id" int8 NOT NULL,
        "is_applied" bool NOT NULL,
        "tstamp" timestamp DEFAULT now() NULL,
        CONSTRAINT "goose_db_version_pkey" PRIMARY KEY (id)
    );

INSERT INTO
    goose_db_version (id, version_id, is_applied, tstamp)
VALUES
    (1, 0, true, now()),
    (2, 1, true, now()),
    (3, 2, true, now()),
    (4, 3, true, now()),
    (5, 4, true, now()),
    (6, 5, true, now()),
    (7, 6, true, now()),
    (8, 7, true, now()) ON CONFLICT
DO NOTHING;