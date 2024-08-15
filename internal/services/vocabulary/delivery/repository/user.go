package repository

import (
	"context"
	"database/sql"
	"fmt"

	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

func (r *VocabRepo) GetVocabulariesByUser(ctx context.Context, userID uuid.UUID) ([]entity.Vocabulary, error) {
	query := `
	SELECT 
		v.id,
		v.user_id,
		v.name,
		v.native_lang,
		v.translate_lang,
		v.description,
		array_agg(tg."text") as tags,
		v.access,
		v.updated_at,
		v.created_at
	FROM vocabulary v
	LEFT JOIN "tag" tg ON tg.id = ANY(v.tags)
	WHERE user_id=$1
	GROUP BY v.id;`
	rows, err := r.pgxPool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vocabularies []entity.Vocabulary
	for rows.Next() {
		var vocab entity.Vocabulary
		var sqlTags []sql.NullString
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.Name,
			&vocab.NativeLang,
			&vocab.TranslateLang,
			&vocab.Description,
			&sqlTags,
			&vocab.Access,
			&vocab.UpdatedAt,
			&vocab.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetVocabularies - scan: %w", err)
		}

		for _, t := range sqlTags {
			if t.Valid {
				vocab.Tags = append(vocab.Tags, entityTag.Tag{Text: t.String})
			}
		}

		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}
