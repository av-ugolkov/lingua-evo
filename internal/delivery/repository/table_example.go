package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type ExampleDB interface {
	AddExample(ctx context.Context, wordId uuid.UUID, example string) (uuid.UUID, error)
}

func (d *Database) AddExample(ctx context.Context, wordId uuid.UUID, example string) (uuid.UUID, error) {
	var id uuid.UUID
	query := `INSERT INTO example (word_id, example) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING id`
	err := d.db.QueryRowContext(ctx, query, wordId, example).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return id, nil
}
