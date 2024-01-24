package repository

import (
	"context"
	"database/sql"
	"fmt"

	entity "lingua-evo/internal/services/tag"

	"github.com/google/uuid"
)

type TagRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *TagRepo {
	return &TagRepo{
		db: db,
	}
}

func (r *TagRepo) AddTag(ctx context.Context, id uuid.UUID, text string) (uuid.UUID, error) {
	query := `
	WITH s AS (
    SELECT id, text FROM tag WHERE text = $2),
	i AS (
    INSERT INTO tag (id, text)
    SELECT $1, $2
    WHERE NOT EXISTS (SELECT 1 FROM s)
    RETURNING id)
	SELECT id
	FROM i
	UNION ALL
		SELECT id
		FROM s;`
	var tid uuid.UUID
	err := r.db.QueryRowContext(ctx, query, id, text).Scan(&tid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("example.repository.TagRepo.AddTag: %w", err)
	}
	return tid, nil
}

func (r *TagRepo) FindTag(ctx context.Context, text string) ([]*entity.Tag, error) {
	query := `SELECT id, text FROM tag WHERE text LIKE '$1%'`
	rows, err := r.db.QueryContext(ctx, query, text)
	if err != nil {
		return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags: %w", err)
	}
	var tags []*entity.Tag
	for rows.Next() {
		var tag *entity.Tag
		err = rows.Scan(&tag.ID, &tag.Text)
		if err != nil {
			return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags - scan: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepo) GetAllTags(ctx context.Context) ([]*entity.Tag, error) {
	query := `SELECT id, text FROM tag`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags: %w", err)
	}
	var tags []*entity.Tag
	for rows.Next() {
		var tag *entity.Tag
		err = rows.Scan(&tag.ID, &tag.Text)
		if err != nil {
			return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags - scan: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
