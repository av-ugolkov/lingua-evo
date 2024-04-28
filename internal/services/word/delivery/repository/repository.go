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

func (r *WordRepo) GetWord(ctx context.Context, vocabID, nativeID uuid.UUID) (entity.VocabWord, error) {
	const query = `SELECT id, translate_ids, example_ids FROM word WHERE vocabulary_id=$1 and native_id=$2;`

	var word entity.VocabWord
	err := r.db.QueryRowContext(ctx, query, vocabID, nativeID).Scan(
		&word.ID,
		pq.Array(&word.TranslateIDs),
		pq.Array(&word.ExampleIDs))
	if err != nil {
		return word, fmt.Errorf("word.repository.WordRepo.GetWord: %w", err)
	}

	return word, nil
}

func (r *WordRepo) AddWord(ctx context.Context, word entity.VocabWord) error {
	const query = `INSERT INTO word (id, vocabulary_id, native_id, translate_ids, example_ids, updated_at, created_at) VALUES($1, $2, $3, $4, $5, $6, $6);`
	_, err := r.db.ExecContext(ctx, query, word.ID, word.VocabID, word.NativeID, word.TranslateIDs, word.ExampleIDs, time.Now().UTC())
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

func (r *WordRepo) GetWordsFromVocabulary(ctx context.Context, dictID uuid.UUID, capacity int) ([]string, error) {
	const query = `
	SELECT text 
	FROM dictionary 
	WHERE id=any(
		SELECT native_id
		FROM word 
		WHERE vocabulary_id=$1
			ORDER BY random() LIMIT $2)`
	rows, err := r.db.QueryContext(ctx, query, dictID, capacity)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetWordsFromVocabulary: %w", err)
	}
	defer rows.Close()

	words := make([]string, 0, capacity)
	for rows.Next() {
		var word string
		err = rows.Scan(&word)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetWordsFromVocabulary - scan: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

func (r *WordRepo) GetRandomWord(ctx context.Context, vocabID uuid.UUID) (entity.VocabWord, error) {
	var vocabWord entity.VocabWord
	query := `SELECT native_id, translate_ids, example_ids FROM word WHERE vocabulary_id=$1 ORDER BY random() LIMIT 1;`
	err := r.db.QueryRowContext(ctx, query, vocabID).Scan(
		&vocabWord.NativeID,
		&vocabWord.TranslateIDs,
		&vocabWord.ExampleIDs,
	)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.repository.WordRepo.GetRandomWord - scan: %w", err)
	}
	return vocabWord, nil
}

func (r *WordRepo) DeleteWord(ctx context.Context, vocabWord entity.VocabWord) error {
	query := `DELETE FROM word WHERE vocabulary_id=$1 AND id=$2;`
	result, err := r.db.ExecContext(ctx, query, vocabWord.VocabID, vocabWord.ID)
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

func (r *WordRepo) GetRandomVocabulary(ctx context.Context, vocabID uuid.UUID, limit int) ([]entity.VocabWord, error) {
	query := `SELECT native_id, translate_ids, example_ids FROM word WHERE vocabulary_id=$1 ORDER BY RANDOM() LIMIT $2;`
	rows, err := r.db.QueryContext(ctx, query, vocabID, limit)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetRandomVocabulary: %w", err)
	}
	defer rows.Close()

	vocabularies := make([]entity.VocabWord, 0, limit)
	for rows.Next() {
		var vocabulary entity.VocabWord
		err = rows.Scan(&vocabulary.NativeID, pq.Array(&vocabulary.TranslateIDs), pq.Array(&vocabulary.ExampleIDs))
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetRandomVocabulary - scan: %w", err)
		}
		vocabularies = append(vocabularies, vocabulary)
	}

	return vocabularies, nil
}

func (r *WordRepo) GetVocabulary(ctx context.Context, vocabID uuid.UUID) ([]entity.VocabWord, error) {
	query := `SELECT id, native_id, translate_ids, example_ids FROM word WHERE vocabulary_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabulary: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	vocabularies := make([]entity.VocabWord, 0, 25)
	for rows.Next() {
		var vocabulary entity.VocabWord
		err = rows.Scan(&vocabulary.ID, &vocabulary.NativeID, pq.Array(&vocabulary.TranslateIDs), pq.Array(&vocabulary.ExampleIDs))
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetVocabulary - scan: %w", err)
		}
		vocabularies = append(vocabularies, vocabulary)
	}

	return vocabularies, nil
}

func (r *WordRepo) UpdateWord(ctx context.Context, vocabWord entity.VocabWord) error {
	query := `UPDATE word SET translate_ids=$1, example_ids=$2 WHERE vocabulary_id=$3 AND native_id=$4;`

	result, err := r.db.ExecContext(ctx, query, vocabWord.TranslateIDs, vocabWord.ExampleIDs, vocabWord.VocabID, vocabWord.NativeID)
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
