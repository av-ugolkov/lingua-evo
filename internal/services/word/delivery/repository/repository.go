package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/word"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

type WordRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *WordRepo {
	return &WordRepo{
		db: db,
	}
}

func (r *WordRepo) GetWord(ctx context.Context, vocabID, nativeID uuid.UUID) (entity.Word, error) {
	const query = `SELECT native_word, translate_word, examples FROM word WHERE vocabulary_id=$1 and native_id=$2;`

	var word entity.Word
	err := r.db.QueryRowContext(ctx, query, vocabID, nativeID).Scan(
		&word.NativeID,
		pq.Array(&word.TranslateWords),
		pq.Array(&word.Examples))
	if err != nil {
		return word, fmt.Errorf("word.repository.WordRepo.GetWord: %w", err)
	}

	return word, nil
}

func (r *WordRepo) AddWord(ctx context.Context, word entity.Word) error {
	const query = `INSERT INTO word (id, vocabulary_id, native_id, translate_words, examples, updated_at, created_at) VALUES($1, $2, $3, $4, $5, $6, $6);`
	_, err := r.db.ExecContext(ctx, query, word.ID, word.VocabID, word.NativeID, word.TranslateWords, word.Examples, time.Now().UTC())
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return fmt.Errorf("word.repository.WordRepo.AddWord: %w", entity.ErrDuplicate)
		default:
			return fmt.Errorf("word.repository.WordRepo.AddWord: %w", err)
		}
	}

	return nil
}

func (r *WordRepo) EditVocabulary(ctx context.Context, vocabulary entity.Word) (int64, error) {
	return 0, nil
}

func (r *WordRepo) GetWordsFromVocabulary(ctx context.Context, dictID uuid.UUID, capacity int) ([]string, error) {
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
		return nil, fmt.Errorf("word.repository.WordRepo.GetWordsFromDictionary: %w", err)
	}
	defer rows.Close()

	words := make([]string, 0, capacity)
	for rows.Next() {
		var word string
		err = rows.Scan(&word)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetWordsFromDictionary - scan: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

func (r *WordRepo) GetRandomWord(ctx context.Context, vocadulary *entity.Word) (*entity.Word, error) {
	query := `SELECT native_word, translate_word, examples FROM vocabulary WHERE dictionary_id=$1 ORDER BY random() LIMIT 1;`
	err := r.db.QueryRowContext(ctx, query, vocadulary.ID).Scan(
		&vocadulary.NativeID,
		&vocadulary.TranslateWords,
		&vocadulary.Examples,
	)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetRandomWord - scan: %w", err)
	}
	return vocadulary, nil
}

func (r *WordRepo) DeleteWord(ctx context.Context, vocabulary entity.Word) error {
	query := `DELETE FROM word WHERE vocabulary_id=$1 AND native_id=$2;`

	result, err := r.db.Exec(query, vocabulary.ID, vocabulary.NativeID)
	if err != nil {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord - exec: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord - rows affected: %w", sql.ErrNoRows)
	} else if err != nil {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord: %w", err)
	}
	return nil
}

func (r *WordRepo) GetRandomVocabulary(ctx context.Context, vocabID uuid.UUID, limit int) ([]entity.Word, error) {
	query := `SELECT native_id, translate_word, examples FROM vocabulary WHERE vocabulary_id=$1 ORDER BY RANDOM() LIMIT $2;`
	rows, err := r.db.QueryContext(ctx, query, vocabID, limit)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetRandomVocabulary: %w", err)
	}
	defer rows.Close()

	vocabularies := make([]entity.Word, 0, limit)
	for rows.Next() {
		var vocabulary entity.Word
		err = rows.Scan(&vocabulary.NativeID, pq.Array(&vocabulary.TranslateWords), pq.Array(&vocabulary.Examples))
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetRandomVocabulary - scan: %w", err)
		}
		vocabularies = append(vocabularies, vocabulary)
	}

	return vocabularies, nil
}

func (r *WordRepo) GetVocabulary(ctx context.Context, vocabID uuid.UUID) ([]entity.Word, error) {
	query := `SELECT native_id, translate_word, examples FROM vocabulary WHERE vocabulary_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabulary: %w", err)
	}
	defer rows.Close()

	vocabularies := make([]entity.Word, 0, 25)
	for rows.Next() {
		var vocabulary entity.Word
		err = rows.Scan(&vocabulary.NativeID, pq.Array(&vocabulary.TranslateWords), pq.Array(&vocabulary.Examples))
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetVocabulary - scan: %w", err)
		}
		vocabularies = append(vocabularies, vocabulary)
	}

	return vocabularies, nil
}

func (r *WordRepo) UpdateWord(ctx context.Context, vocabulary entity.Word) error {
	query := `UPDATE vocabulary SET translate_word=$1, examples=$2 WHERE vocabulary_id=$4 AND native_id=$5;`

	result, err := r.db.ExecContext(ctx, query, vocabulary.TranslateWords, vocabulary.Examples, vocabulary.ID, vocabulary.NativeID)
	if err != nil {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord - exec: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord - rows affected: %w", sql.ErrNoRows)
	} else if err != nil {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord: %w", err)
	}
	return nil
}
