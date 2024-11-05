CREATE TABLE IF NOT EXISTS
    vocabulary_notifications (
        user_id UUID REFERENCES users (id) ON DELETE CASCADE,
        vocab_id UUID REFERENCES vocabulary (id) ON DELETE CASCADE,
        created_at timestamp NOT NULL
    );

CREATE UNIQUE INDEX IF NOT EXISTS "idx_unique_vocabulary_notifications__user_id_vocab_id" ON "vocabulary_notifications" ("user_id", "vocab_id");