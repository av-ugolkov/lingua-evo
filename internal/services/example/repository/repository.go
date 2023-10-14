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

func (r *ExampleRepo) AddExample(ctx context.Context, id uuid.UUID, text, langCode string) error {
	query := fmt.Sprintf(`INSERT INTO example_%s (id, text) VALUES($1, $2) ON CONFLICT DO NOTHING`, langCode)
	err := r.db.QueryRowContext(ctx, query, id, text).Scan(&id)
	if err != nil {
		return fmt.Errorf("example.repository.ExampleRepo.AddExample: %w", err)
	}

	return nil
}

func (r *ExampleRepo) GetExample(ctx context.Context, id uuid.UUID, langCode string) (string, error) {
	var text string
	query := fmt.Sprintf(`SELECT text FROM example_%s WHERE id=$1`, langCode)
	err := r.db.QueryRowContext(ctx, query, id).Scan(&text)
	if err != nil {
		return "", fmt.Errorf("example.repository.ExampleRepo.AddExample: %w", err)
	}

	return text, nil
}
