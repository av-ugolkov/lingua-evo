package repository

import (
	"context"
	"database/sql"

	entityWord "lingua-evo/internal/services/word/entity"

	"github.com/google/uuid"
)

type VocabularyRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *VocabularyRepo {
	return &VocabularyRepo{
		db: db,
	}
}

func (r *VocabularyRepo) AddWord(ctx context.Context, dictID, originalWord uuid.UUID, translateWord []uuid.UUID, examples []uuid.UUID, tags []uuid.UUID) error {
	const query = `INSERT INTO vocabulary (dictionary_id, original_word, translate_word, examples, tags) VALUES($1, $2, $3, $4, $5)`
	_, err := r.db.QueryContext(ctx, query, dictID, originalWord, translateWord, examples, tags)
	if err != nil {
		return err
	}

	return nil
}

func (r *VocabularyRepo) GetWordsByUser(ctx context.Context, userId string) ([]entityWord.Word, error) {
	return []entityWord.Word{}, nil
}

func (r *VocabularyRepo) GetRandomWordByUser(ctx context.Context, userId string) (entityWord.Word, error) {
	return entityWord.Word{}, nil
}
