package repository

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

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

func (r *VocabularyRepo) AddWord(ctx context.Context, vocabulary entity.Vocabulary) error {
	const query = `INSERT INTO vocabulary (dictionary_id, native_word, translate_word, examples, tags) VALUES($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING;`
	rows, err := r.db.QueryContext(ctx, query, vocabulary.DictionaryId, vocabulary.NativeWord, vocabulary.TranslateWords, vocabulary.Examples, vocabulary.Tags)
	if err != nil {
		return fmt.Errorf("vocabulary.repository.VocabularyRepo.AddWord: %w", err)
	}
	defer rows.Close()

	return nil
}

func (r *VocabularyRepo) EditVocabulary(ctx context.Context, vocabulary entity.Vocabulary) (int64, error) {
	return 0, nil
}

func (r *VocabularyRepo) GetWordsFromDictionary(ctx context.Context, dictID uuid.UUID, capacity int) ([]string, error) {
	const query = `
	SELECT text 
	FROM word 
	WHERE id=any(
		SELECT native_word 
		FROM vocabulary 
		WHERE dictionary_id=$1
			ORDER BY random() LIMIT $2)`
	rows, err := r.db.QueryContext(ctx, query, dictID, capacity)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabularyRepo.GetWordsFromDictionary: %w", err)
	}
	defer rows.Close()

	words := make([]string, 0, capacity)
	for rows.Next() {
		var word string
		err = rows.Scan(&word)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabularyRepo.GetWordsFromDictionary - scan: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
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

func (r *VocabularyRepo) DeleteWord(ctx context.Context, vocabulary entity.Vocabulary) error {
	query := `DELETE FROM vocabulary WHERE dictionary_id=$1 AND native_word=$2;`

	result, err := r.db.Exec(query, vocabulary.DictionaryId, vocabulary.NativeWord)
	if err != nil {
		return fmt.Errorf("vocabulary.repository.VocabularyRepo.DeleteWord - exec: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("vocabulary.repository.VocabularyRepo.DeleteWord - rows affected: %w", sql.ErrNoRows)
	} else if err != nil {
		return fmt.Errorf("vocabulary.repository.VocabularyRepo.DeleteWord: %w", err)
	}
	return nil
}

func (r *VocabularyRepo) GetWords(ctx context.Context, dictID uuid.UUID, limit int) ([]entity.Vocabulary, error) {
	query := `SELECT native_word, translate_word, examples, tags FROM vocabulary WHERE dictionary_id=$1 ORDER BY RANDOM() LIMIT $2;`
	rows, err := r.db.QueryContext(ctx, query, dictID, limit)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabularyRepo.GetWords: %w", err)
	}
	defer rows.Close()

	words := make([]entity.Vocabulary, 0, limit)
	for rows.Next() {
		var word entity.Vocabulary
		err = rows.Scan(&word.NativeWord, &word.TranslateWords, &word.Examples, &word.Tags)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabularyRepo.GetWords - scan: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

func (r *VocabularyRepo) UpdateWord(ctx context.Context, vocabulary entity.Vocabulary) error {
	query := `UPDATE vocabulary SET translate_word=$1, examples=$2, tags=$3 WHERE dictionary_id=$4 AND native_word=$5;`

	result, err := r.db.ExecContext(ctx, query, vocabulary.TranslateWords, vocabulary.Examples, vocabulary.Tags, vocabulary.DictionaryId, vocabulary.NativeWord)
	if err != nil {
		return fmt.Errorf("vocabulary.repository.VocabularyRepo.UpdateWord - exec: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("vocabulary.repository.VocabularyRepo.UpdateWord - rows affected: %w", sql.ErrNoRows)
	} else if err != nil {
		return fmt.Errorf("vocabulary.repository.VocabularyRepo.UpdateWord: %w", err)
	}
	return nil
}
