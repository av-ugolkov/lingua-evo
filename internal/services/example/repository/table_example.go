package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ExampleRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *ExampleRepo {
	return &ExampleRepo{
		db: db,
	}
}

func (r *ExampleRepo) AddExample(ctx context.Context, wordId uuid.UUID, example string) (uuid.UUID, error) {
	var id uuid.UUID
	query := `INSERT INTO example (word_id, example) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING id`
	err := r.db.QueryRowContext(ctx, query, wordId, example).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return id, nil
}
