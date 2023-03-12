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
	query := `insert into users (user_id, user_mame) values ($1, $2) ON CONFLICT DO NOTHING`
	_, err := d.db.QueryContext(ctx, query, userId, userName)
	if err != nil {
		return fmt.Errorf("database.AddUser.QueryRow: %w", err)
	}
	return nil
}
