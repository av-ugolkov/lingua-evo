-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_en__id" ON "dictionary_en" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_ru__id" ON "dictionary_ru" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_fi__id" ON "dictionary_fi" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_fr__id" ON "dictionary_fr" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_it__id" ON "dictionary_it" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_es__id" ON "dictionary_es" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_de__id" ON "dictionary_de" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_pt__id" ON "dictionary_pt" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_dictionary_sv__id" ON "dictionary_sv" ("id");

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_en__id" ON "example_en" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_ru__id" ON "example_ru" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_fi__id" ON "example_fi" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_fr__id" ON "example_fr" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_es__id" ON "example_es" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_it__id" ON "example_it" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_de__id" ON "example_de" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_pt__id" ON "example_pt" ("id");
CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_example_sv__id" ON "example_sv" ("id");

-- +goose Down
DROP INDEX "idx_unique_dictionary_en__id";
DROP INDEX "idx_unique_dictionary_ru__id";
DROP INDEX "idx_unique_dictionary_fi__id";
DROP INDEX "idx_unique_dictionary_fr__id";
DROP INDEX "idx_unique_dictionary_it__id";
DROP INDEX "idx_unique_dictionary_es__id";
DROP INDEX "idx_unique_dictionary_de__id";
DROP INDEX "idx_unique_dictionary_pt__id";
DROP INDEX "idx_unique_dictionary_sv__id";

DROP INDEX "idx_unique_example_en__id";
DROP INDEX "idx_unique_example_ru__id";
DROP INDEX "idx_unique_example_fi__id";
DROP INDEX "idx_unique_example_fr__id";
DROP INDEX "idx_unique_example_it__id";
DROP INDEX "idx_unique_example_es__id";
DROP INDEX "idx_unique_example_de__id";
DROP INDEX "idx_unique_example_pt__id";
DROP INDEX "idx_unique_example_sv__id";
