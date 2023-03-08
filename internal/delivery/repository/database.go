package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	db *pgxpool.Pool
}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{
		db: pool,
	}
}

func (d *Database) AddUser(ctx context.Context, userId int, userName string) error {
	query := `insert into users (user_id, user_mame) values ($1, $2) ON CONFLICT DO NOTHING`
	_, err := d.db.Query(ctx, query, userId, userName)
	if err != nil {
		return fmt.Errorf("database.AddUser.QueryRow: %w", err)
	}
	return nil
}
