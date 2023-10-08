package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lingua-evo/internal/services/word/entity"
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
		return uuid.Nil, fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return id, nil
}

func (r *WordRepo) GetWord(ctx context.Context, text, langCode string) (*entity.Word, error) {
	word := &entity.Word{}
	table := getTable(langCode)
	query := fmt.Sprintf(`SELECT id, text, pronunciation, lang_code FROM "%s" WHERE text=$1 AND lang_code=$2;`, table)
	err := r.db.QueryRowContext(ctx, query, text, langCode).Scan(&word.ID, &word.Text, &word.Pronunciation, &word.LanguageCode)
	if err != nil {
		return nil, err
	}
	return word, nil
}

func (r *WordRepo) EditWord(ctx context.Context, w *entity.Word) error {
	return nil
}

func (r *WordRepo) FindWord(ctx context.Context, w *entity.Word) (uuid.UUID, error) {
	var id uuid.UUID
	query := `SELECT id FROM word WHERE text=$1 AND lang_code=$2`
	err := r.db.QueryRowContext(ctx, query, w.Text, w.LanguageCode).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("database.FindWord.QueryRow: %w", err)
	}
	return id, nil
}

func (r *WordRepo) FindWords(ctx context.Context, w string) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	query := `SELECT id FROM word WHERE text=$1`
	err := r.db.QueryRowContext(ctx, query, w).Scan(ids)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.FindWord - scan: %w", err)
	}

	return ids, nil
}

func (r *WordRepo) RemoveWord(ctx context.Context, w *entity.Word) error {
	return nil
}

func (r *WordRepo) GetRandomWord(ctx context.Context, lang string) (*entity.Word, error) {
	table := getTable(lang)
	query := fmt.Sprintf(`SELECT text FROM "%s" ORDER BY RANDOM() LIMIT 1;`, table)
	word := &entity.Word{}
	err := r.db.QueryRowContext(ctx, query).Scan(&word.Text)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetRandomWord - scan: %w", err)
	}
	return word, nil
}

func (r *WordRepo) SharedWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	return nil, nil
}

func getTable(langCode string) string {
	table := "word"
	if len(langCode) != 0 {
		table = fmt.Sprintf(`%s_%s`, table, langCode)
	}
	return table
}
