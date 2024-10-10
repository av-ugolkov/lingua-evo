package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(cfg *pgxpool.Config) (*pgxpool.Pool, error) {
	connPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("create DB connection pool: %w", err)
	}
	return connPool, nil
}
