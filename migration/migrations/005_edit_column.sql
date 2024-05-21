-- +goose Up
ALTER TABLE "word"
ADD COLUMN "pronunciation" TEXT;

-- +goose Down
ALTER TABLE "word"
DROP COLUMN "pronunciation";