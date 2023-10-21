package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lingua-evo/internal/services/lingua/vocabulary/entity"
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
	const query = `INSERT INTO vocabulary (dictionary_id, native_word, translate_word, examples, tags) VALUES($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING;`
	_, err := r.db.QueryContext(ctx, query, vocabulary.DictionaryId, vocabulary.NativeWord, vocabulary.TranslateWords, vocabulary.Examples, vocabulary.Tags)
	if err != nil {
		return err
	}

	return nil
}

func (r *VocabularyRepo) EditVocabulary(ctx context.Context, vocabulary entity.Vocabulary) (int64, error) {
	return 0, nil
}

func (r *VocabularyRepo) GetWords(ctx context.Context, dictID string) ([]*entity.Vocabulary, error) {
	return []*entity.Vocabulary{}, nil
}

func (r *VocabularyRepo) GetRandomWord(ctx context.Context, vocadulary *entity.Vocabulary) (*entity.Vocabulary, error) {
	query := `SELECT * FROM vocabulary WHERE dictionary_id=$1 ORDER BY random() LIMIT 1;`
	err := r.db.QueryRowContext(ctx, query, vocadulary.DictionaryId).Scan(
		&vocadulary.NativeWord,
		&vocadulary.TranslateWords,
		&vocadulary.Examples,
		&vocadulary.Tags)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabularyRepo.GetRandomWord - scan: %w", err)
	}
	return vocadulary, nil
}

func (r *VocabularyRepo) DeleteWord(ctx context.Context, vocabulary entity.Vocabulary) (int64, error) {
	query := `DELETE FROM vocabulary WHERE dictopnary_id=$1 AND native_word=$2;`

	result, err := r.db.Exec(query, vocabulary.DictionaryId, vocabulary.NativeWord)
	if err != nil {
		return 0, fmt.Errorf("vocabulary.repository.VocabularyRepo.DeleteWord - exec: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("vocabulary.repository.VocabularyRepo.DeleteWord - rows affected: %w", err)
	}
	return rows, nil
}
