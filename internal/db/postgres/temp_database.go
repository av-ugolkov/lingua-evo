package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TempPostgres struct {
	container *postgres.PostgresContainer
	PgxPool   *pgxpool.Pool
}

var (
	Instance *TempPostgres
	mx       sync.Mutex
	wg       sync.WaitGroup
)

func NewTempPostgres(ctx context.Context, root string) *TempPostgres {
	mx.Lock()
	defer mx.Unlock()
	wg.Add(1)

	var err error
	if Instance != nil {
		return Instance
	}

	Instance = &TempPostgres{}
	Instance.container, err = postgres.Run(ctx,
		"docker.io/postgres:16.1-alpine3.19",
		postgres.WithInitScripts(
			path.Join(root, "migration/migrations/001_create_tables.up.sql"),
			path.Join(root, "migration/migrations/002_add_data.up.sql"),
			//path.Join(root, "migration/migrations/003_insert_dictionary_words.up.sql"),
			path.Join(root, "migration/migrations/004_unique_index.up.sql"),
			path.Join(root, "migration/migrations/005_edit_column.up.sql"),
			path.Join(root, "migration/migrations/006_access_vocabulary.up.sql"),
			path.Join(root, "migration/migrations/007_drop_goose.up.sql"),
			path.Join(root, "migration/migrations/008_notifications.up.sql"),
			path.Join(root, "migration/migrations/009_word_desc.up.sql"),
			path.Join(root, "migration/migrations/010_event_vocab.up.sql")),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second)),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to start container: %s", err))
		return nil
	}

	connStr, err := Instance.container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get connection string: %s", err))
		return nil
	}

	Instance.PgxPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create pgxpool: %s", err))
		return nil
	}

	return Instance
}

func (t *TempPostgres) DropDB(ctx context.Context) {
	go func() {
		wg.Done()
		time.Sleep(1 * time.Second)

		wg.Wait()
		if err := t.container.Terminate(ctx); err != nil {
			slog.Error(fmt.Sprintf("failed to terminate container: %s", err))
		}

		Instance = nil
	}()
}

func (t *TempPostgres) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	rows, err := t.PgxPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return rows, nil
}
