package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/word"
)

type WordRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *WordRepo {
	return &WordRepo{
		db: db,
	}
}

func (r *WordRepo) AddWord(ctx context.Context, w *entity.Word) (uuid.UUID, error) {
	var id uuid.UUID
	table := getTable(w.LanguageCode)
	query := fmt.Sprintf(`INSERT INTO "%s" (id, text, pronunciation, lang_code, created_at) VALUES($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING RETURNING id;`, table)
	err := r.db.QueryRowContext(ctx, query, w.ID, w.Text, w.Pronunciation, w.LanguageCode, time.Now().UTC()).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("word.repository.WordRepo.AddWord - query: %w", err)
	}

	return id, nil
}

func (r *WordRepo) GetWordByText(ctx context.Context, w *entity.Word) (uuid.UUID, error) {
	word := &entity.Word{}
	table := getTable(w.LanguageCode)
	query := fmt.Sprintf(`SELECT id FROM "%s" WHERE text=$1 AND lang_code=$2;`, table)
	err := r.db.QueryRowContext(ctx, query, w.Text, w.LanguageCode).Scan(&word.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("word.repository.WordRepo.GetWordByText - query: %w", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil
	}
	return word.ID, nil
}

func (r *WordRepo) GetWords(ctx context.Context, ids []uuid.UUID) ([]entity.Word, error) {
	query := `SELECT id, text, pronunciation FROM word WHERE id=ANY($1);`
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetWords - query: %w", err)
	}
	defer rows.Close()

	words := make([]entity.Word, 0, len(ids))
	for rows.Next() {
		var word entity.Word
		err = rows.Scan(&word.ID, &word.Text, &word.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetWords - scan: %w", err)
		}
		words = append(words, word)
	}
	return words, nil
}

func (r *WordRepo) UpdateWord(ctx context.Context, w *entity.Word) error {
	query := `UPDATE word SET text=$1, pronunciation=$2 WHERE id=$3`
	result, err := r.db.ExecContext(ctx, query, w.Text, w.Pronunciation, w.ID)
	if err != nil {
		return fmt.Errorf("word.repository.WordRepo.EditWord - exec: %w", err)
	}

	if rows, err := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("word.repository.WordRepo.EditWord: not fount effected rows")
	} else if err != nil {
		return fmt.Errorf("word.repository.WordRepo.EditWord - rows affected: %w", err)
	}

	return nil
}

func (r *WordRepo) FindWords(ctx context.Context, w *entity.Word) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	query := `SELECT id FROM word WHERE text=$1% AND lang_code=$2;`
	rows, err := r.db.QueryContext(ctx, query, w.Text, w.LanguageCode)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.FindWords - query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.FindWords - scan: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *WordRepo) DeleteWord(ctx context.Context, w *entity.Word) (int64, error) {
	query := `DELETE FROM word WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, w.ID)
	if err != nil {
		return 0, fmt.Errorf("word.repository.WordRepo.DeleteWord - exec: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("word.repository.WordRepo.DeleteWord - rows affected: %w", err)
	}
	return rows, nil
}

func (r *WordRepo) GetRandomWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	table := getTable(w.LanguageCode)
	query := fmt.Sprintf(`SELECT text, pronunciation FROM "%s" ORDER BY RANDOM() LIMIT 1;`, table)
	word := &entity.Word{}
	err := r.db.QueryRowContext(ctx, query).Scan(&word.Text, &word.Pronunciation)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetRandomWord - scan: %w", err)
	}

	word.LanguageCode = w.LanguageCode

	return word, nil
}

func (r *WordRepo) SharedWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	// TODO implement me later
	return nil, nil
}

func getTable(langCode string) string {
	table := "word"
	if len(langCode) != 0 {
		table = fmt.Sprintf(`%s_%s`, table, langCode)
	}
	return table
}
