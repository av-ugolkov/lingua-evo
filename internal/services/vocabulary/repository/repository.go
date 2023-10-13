package repository

import (
	"context"
	"database/sql"

	"lingua-evo/internal/services/vocabulary/entity"
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

func (r *VocabularyRepo) AddWord(ctx context.Context, vocabulary entity.Vocabulary) (uuid.UUID, error) {
	var id uuid.UUID
	const query = `INSERT INTO vocabulary (dictionary_id, native_word, translate_word, examples, tags) VALUES($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING id;`
	err := r.db.QueryRowContext(ctx, query, vocabulary.DictionaryId, vocabulary.NativeWord, vocabulary.TranslateWord, vocabulary.Examples, vocabulary.Tags).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *VocabularyRepo) GetWord(ctx context.Context, dictID, word string) (*entityWord.Word, error) {
	return &entityWord.Word{}, nil
}

func (r *VocabularyRepo) GetWords(ctx context.Context, dictID string) ([]*entityWord.Word, error) {
	return []*entityWord.Word{}, nil
}

func (r *VocabularyRepo) GetRandomWord(ctx context.Context, userId string) (*entityWord.Word, error) {
	return &entityWord.Word{}, nil
}
