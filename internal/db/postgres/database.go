package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(cfg config.DbSQL) (*pgxpool.Pool, error) {
	dbConfig, err := pgxpool.ParseConfig(cfg.GetConnStr())
	if err != nil {
		return nil, fmt.Errorf("parse DB connection string: %w", err)
	}

	dbConfig.MaxConns = int32(cfg.MaxConns)
	dbConfig.MinConns = int32(cfg.MinConns)
	dbConfig.MaxConnLifetime = time.Duration(cfg.MaxConnLifetime) * time.Second
	dbConfig.MaxConnIdleTime = time.Duration(cfg.MaxConnIdleTime) * time.Second
	dbConfig.HealthCheckPeriod = time.Duration(cfg.HealthCheckPeriod) * time.Second
	dbConfig.ConnConfig.ConnectTimeout = time.Duration(cfg.ConnectTimeout) * time.Second

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("create DB connection pool: %w", err)
	}

	return connPool, nil
}
