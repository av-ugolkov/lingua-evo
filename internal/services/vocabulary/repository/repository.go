package repository

import (
	"context"
	"database/sql"

	"lingua-evo/internal/services/vocabulary/entity"
	entityWord "lingua-evo/internal/services/word/entity"
)

type VocabularyRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *VocabularyRepo {
	return &VocabularyRepo{
		db: db,
	}
}

func (r *VocabularyRepo) AddWord(ctx context.Context, vocabulary entity.Vocabulary) error {
	const query = `INSERT INTO vocabulary (dictionary_id, original_word, translate_word, examples, tags) VALUES($1, $2, $3, $4, $5)`
	_, err := r.db.QueryContext(ctx, query, vocabulary.DictionaryId, vocabulary.OriginalWord, vocabulary.TranslateWord, vocabulary.Examples, vocabulary.Tags)
	if err != nil {
		return err
	}

	return nil
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
