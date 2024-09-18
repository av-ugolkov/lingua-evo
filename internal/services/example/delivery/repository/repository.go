package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/services/example"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/example"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExampleRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *ExampleRepo {
	return &ExampleRepo{
		pgxPool: pgxPool,
	}
}

func (r *ExampleRepo) AddExamples(ctx context.Context, examples []entity.Example, langCode string) ([]uuid.UUID, error) {
	wordTexts := make([]string, 0, len(examples))
	statements := make([]string, 0, len(examples))
	params := make([]any, 0, len(examples)+1)
	params = append(params, &wordTexts)
	counter := len(params)
	for i := 0; i < len(examples); i++ {
		wordTexts = append(wordTexts, examples[i].Text)
		statement := "$" + strconv.Itoa(counter+1) +
			",$" + strconv.Itoa(counter+2) +
			",$" + strconv.Itoa(counter+3)

		counter += 3
		statements = append(statements, "("+statement+")")

		params = append(params, examples[i].ID, examples[i].Text, examples[i].CreatedAt.Format(time.RFC3339))
	}

	query := fmt.Sprintf(`
	WITH s AS (
    		SELECT id FROM example_%[1]s WHERE text = any($1::text[])),
		ins AS (
    		INSERT INTO example_%[1]s (id, text, created_at)
			VALUES %[2]s
    		ON CONFLICT DO NOTHING RETURNING id)
		SELECT id
		FROM ins
		UNION ALL
		SELECT id
		FROM s;`, langCode, strings.Join(statements, ", "))
	rows, err := r.pgxPool.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("example.repository.ExampleRepo.AddExamples - query: %w", err)
	}
	defer rows.Close()

	tagIDs := make([]uuid.UUID, 0, len(examples))
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("example.repository.ExampleRepo.AddExamples - scan: %w", err)
		}
		tagIDs = append(tagIDs, id)
	}

	return tagIDs, nil
}

func (r *ExampleRepo) GetExampleByValue(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf(`SELECT id FROM example_%s WHERE text=$1`, langCode)
	err := r.pgxPool.QueryRow(ctx, query, text).Scan(&id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("example.repository.ExampleRepo.GetExample: %w", err)
	} else if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, nil
	}

	return id, nil
}

func (r *ExampleRepo) GetExampleById(ctx context.Context, id uuid.UUID, langCode string) (string, error) {
	var text string
	query := fmt.Sprintf(`SELECT text FROM example_%s WHERE id=$1`, langCode)
	err := r.pgxPool.QueryRow(ctx, query, id).Scan(&text)
	if err != nil {
		return "", fmt.Errorf("example.repository.ExampleRepo.AddExample: %w", err)
	}

	return text, nil
}

func (r *ExampleRepo) GetExamples(ctx context.Context, ids []uuid.UUID) ([]example.Example, error) {
	query := `SELECT id, text FROM example WHERE id = ANY($1)`
	rows, err := r.pgxPool.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("example.repository.ExampleRepo.GetExamples: %w", err)
	}
	defer rows.Close()

	examples := make([]example.Example, 0, len(ids))
	for rows.Next() {
		var ex example.Example
		err = rows.Scan(&ex.ID, &ex.Text)
		if err != nil {
			return nil, fmt.Errorf("example.repository.ExampleRepo.GetExamples - scan: %w", err)
		}
		examples = append(examples, ex)
	}

	return examples, nil
}
