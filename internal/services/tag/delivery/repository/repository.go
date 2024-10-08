package repository

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/tag"

	"github.com/google/uuid"
)

type TagRepo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *TagRepo {
	return &TagRepo{
		tr: tr,
	}
}

func (r *TagRepo) AddTag(ctx context.Context, text string) (uuid.UUID, error) {
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
	err := r.tr.QueryRow(ctx, query, uuid.New(), text).Scan(&tid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("example.repository.TagRepo.AddTag: %w", err)
	}
	return tid, nil
}

func (r *TagRepo) FindTag(ctx context.Context, text string) ([]entity.Tag, error) {
	query := `SELECT id, text FROM tag WHERE text LIKE '$1%'`
	rows, err := r.tr.Query(ctx, query, text)
	if err != nil {
		return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags: %w", err)
	}
	defer rows.Close()

	var tags []entity.Tag
	for rows.Next() {
		var tag entity.Tag
		err = rows.Scan(&tag.ID, &tag.Text)
		if err != nil {
			return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags - scan: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepo) GetTag(ctx context.Context, text string) (uuid.UUID, error) {
	query := `SELECT id FROM tag WHERE text = $1`
	var id uuid.UUID
	err := r.tr.QueryRow(ctx, query, text).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("example.repository.TagRepo.GetTag: %w", err)
	}
	return id, nil
}

func (r *TagRepo) GetTagsInVocabulary(ctx context.Context, vocabID uuid.UUID) ([]entity.Tag, error) {
	query := `SELECT id,text FROM tag WHERE id=ANY((SELECT tags FROM vocabulary WHERE id=$1)::uuid[]);`
	rows, err := r.tr.Query(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("example.repository.TagRepo.GetTags: %w", err)
	}
	defer rows.Close()

	tags := make([]entity.Tag, 0)
	for rows.Next() {
		var tag entity.Tag
		err = rows.Scan(&tag.ID, &tag.Text)
		if err != nil {
			return nil, fmt.Errorf("example.repository.TagRepo.GetTags - scan: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *TagRepo) GetAllTags(ctx context.Context) ([]entity.Tag, error) {
	query := `SELECT id, text FROM tag`
	rows, err := r.tr.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags: %w", err)
	}
	defer rows.Close()

	var tags []entity.Tag
	for rows.Next() {
		var tag entity.Tag
		err = rows.Scan(&tag.ID, &tag.Text)
		if err != nil {
			return nil, fmt.Errorf("example.repository.TagRepo.GetAllTags - scan: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
