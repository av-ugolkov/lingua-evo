TRUNCATE TABLE "users" RESTART IDENTITY CASCADE;

TRUNCATE TABLE "language" RESTART IDENTITY CASCADE;

TRUNCATE TABLE "word" RESTART IDENTITY CASCADE;


DROP INDEX IF EXISTS "idx_unique_dictionary_en__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_ru__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_fi__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_fr__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_it__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_es__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_de__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_pt__text";

DROP INDEX IF EXISTS "idx_unique_dictionary_sv__text";

DROP INDEX IF EXISTS "idx_unique_example_en__text";

DROP INDEX IF EXISTS "idx_unique_example_ru__text";

DROP INDEX IF EXISTS "idx_unique_example_fi__text";

DROP INDEX IF EXISTS "idx_unique_example_fr__text";

DROP INDEX IF EXISTS "idx_unique_example_it__text";

DROP INDEX IF EXISTS "idx_unique_example_es__text";

DROP INDEX IF EXISTS "idx_unique_example_de__text";

DROP INDEX IF EXISTS "idx_unique_example_pt__text";

DROP INDEX IF EXISTS "idx_unique_example_sv__text";