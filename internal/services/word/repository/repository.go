package repository

import (
	"context"
	"database/sql"
	"fmt"

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
	query := `INSERT INTO word (text, lang) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING id`
	err := r.db.QueryRowContext(ctx, query, w.Text, w.Language).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return id, nil
}

func (r *WordRepo) EditWord(ctx context.Context, w *entity.Word) error {
	return nil
}

func (r *WordRepo) FindWord(ctx context.Context, w *entity.Word) (uuid.UUID, error) {
	var id uuid.UUID
	query := `SELECT id FROM word WHERE text=$1 AND lang=$2`
	err := r.db.QueryRowContext(ctx, query, w.Text, w.Language).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("database.FindWord.QueryRow: %w", err)
	}
	return id, nil
}

func (r *WordRepo) FindWords(ctx context.Context, w string) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	query := `SELECT id FROM word WHERE text=$1`
	err := r.db.QueryRowContext(ctx, query, w).Scan(&ids)
	if err != nil {
		return []uuid.UUID{}, fmt.Errorf("database.FindWord.QueryRow: %w", err)
	}
	return ids, nil
}

func (r *WordRepo) RemoveWord(ctx context.Context, w *entity.Word) error {
	return nil
}

func (r *WordRepo) PickRandomWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	return nil, nil
}

func (r *WordRepo) SharedWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	return nil, nil
}
