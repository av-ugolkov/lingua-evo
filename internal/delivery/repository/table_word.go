package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type WordDB interface {
	AddWord(ctx context.Context, w *Word) (uuid.UUID, error)
	EditWord(ctx context.Context, w *Word) error
	FindWord(ctx context.Context, w string) (*Word, error)
	RemoveWord(ctx context.Context, w *Word) error
	PickRandomWord(ctx context.Context, w *Word) (*Word, error)
	SharedWord(ctx context.Context, w *Word) (*Word, error)
}

func (d *Database) AddWord(ctx context.Context, w *Word) (uuid.UUID, error) {
	query := `INSERT INTO word (text, lang) VALUES($1, $2) RETURNING id`
	var id uuid.UUID
	err := d.db.QueryRowContext(ctx, query, w.Text, w.Language).Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return id, nil
}

func (d *Database) EditWord(ctx context.Context, w *Word) error {
	return nil
}

func (d *Database) FindWord(ctx context.Context, w string) (*Word, error) {
	return nil, nil
}

func (d *Database) RemoveWord(ctx context.Context, w *Word) error {
	return nil
}

func (d *Database) PickRandomWord(ctx context.Context, w *Word) (*Word, error) {
	return nil, nil
}

func (d *Database) SharedWord(ctx context.Context, w *Word) (*Word, error) {
	return nil, nil
}
