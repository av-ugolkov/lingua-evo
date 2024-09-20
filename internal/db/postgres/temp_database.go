package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	// "github.com/testcontainers/testcontainers-go"
	// "github.com/testcontainers/testcontainers-go/modules/postgres"
	// "github.com/testcontainers/testcontainers-go/wait"
)

type TempPostgres struct {
	// container *postgres.PostgresContainer
	pgxPool *pgxpool.Pool
}

func NewTempPostgres(ctx context.Context, dbName string) *TempPostgres {
	// dbUser := "user"
	// dbPassword := "password"

	// var err error
	tempP := TempPostgres{}
	// tempP.container, err = postgres.Run(ctx,
	// 	"docker.io/postgres:16.1-alpine3.19",
	// 	postgres.WithDatabase(dbName),
	// 	postgres.WithUsername(dbUser),
	// 	postgres.WithPassword(dbPassword),
	// 	testcontainers.WithWaitStrategy(
	// 		wait.ForLog("database system is ready to accept connections").
	// 			WithOccurrence(2).
	// 			WithStartupTimeout(15*time.Second)),
	// )
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("failed to start container: %s", err))
	// 	return nil
	// }

	// connStr, err := tempP.container.ConnectionString(ctx, "sslmode=disable")
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("failed to get connection string: %s", err))
	// 	return nil
	// }

	// tempP.pgxPool, err = pgxpool.New(ctx, connStr)
	// if err != nil {
	// 	slog.Error(fmt.Sprintf("failed to create pgxpool: %s", err))
	// 	return nil
	// }

	return &tempP
}

func (t *TempPostgres) DropDB(ctx context.Context) {
	// if err := t.container.Terminate(ctx); err != nil {
	// 	slog.Error(fmt.Sprintf("failed to terminate container: %s", err))
	// }
}

func (t *TempPostgres) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return t.pgxPool.Exec(ctx, query, args...)
}
