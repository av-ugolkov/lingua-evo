package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

func (d *Database) AddUser(ctx context.Context, userId int, userName string) error {
	query := `INSERT INTO users (user_id, user_mame) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	err := d.db.QueryRowContext(ctx, query, userId, userName)
	if err != nil {
		return fmt.Errorf("database.AddUser.QueryRow: %w", err)
	}
	return nil
}
