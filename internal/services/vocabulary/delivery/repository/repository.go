package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type VocabRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *VocabRepo {
	return &VocabRepo{
		pgxPool: pgxPool,
	}
}

func (r *VocabRepo) Add(ctx context.Context, vocab entity.Vocabulary, tagIDs []uuid.UUID) error {
	query := `INSERT INTO vocabulary (id, user_id, name, native_lang, translate_lang, tags, updated_at, created_at, access, access_edit) VALUES($1, $2, $3, $4, $5, $6, $7, $7, $8, $9)`

	_, err := r.pgxPool.Exec(ctx, query, vocab.ID, vocab.UserID, vocab.Name, vocab.NativeLang, vocab.TranslateLang, pq.Array(tagIDs), time.Now().UTC(), vocab.Access, false)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Add: %w", err)
	}

	return nil
}

func (r *VocabRepo) Delete(ctx context.Context, vocab entity.Vocabulary) error {
	query := `DELETE FROM vocabulary WHERE user_id=$1 AND name=$2;`
	result, err := r.pgxPool.Exec(ctx, query, vocab.UserID, vocab.Name)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Delete: %w", err)
	}
	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Delete: %w", entity.ErrVocabularyNotFound)
	}

	return nil
}

func (r *VocabRepo) Get(ctx context.Context, vocabID uuid.UUID) (entity.Vocabulary, error) {
	query := `SELECT id, user_id, name, n.lang as native_lang, t.lang as translate_lang FROM vocabulary v
left join "language" n on n.code=v.native_lang 
left join "language" t on t.code=v.translate_lang 
WHERE id=$1;`
	var vocab entity.Vocabulary
	err := r.pgxPool.QueryRow(ctx, query, vocabID).Scan(&vocab.ID, &vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Get: %w", err)
	}
	return vocab, nil
}

func (r *VocabRepo) GetByName(ctx context.Context, userID uuid.UUID, name string) (entity.Vocabulary, error) {
	query := `SELECT id, user_id, name, n.lang as native_lang, t.lang as translate_lang FROM vocabulary v
left join "language" n on n.code =v.native_lang 
left join "language" t on t.code =v.translate_lang 
WHERE user_id=$1 AND name=$2;`
	var vocab entity.Vocabulary
	err := r.pgxPool.QueryRow(ctx, query, userID, name).Scan(&vocab.ID, &vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
	}
	return vocab, nil
}

func (r *VocabRepo) GetTagsVocabulary(ctx context.Context, vocabID uuid.UUID) ([]string, error) {
	query := `SELECT "text" from tag t
left join vocabulary v on t.id=any(v.tags)
where v.id=$1;`
	rows, err := r.pgxPool.Query(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *VocabRepo) GetByID(ctx context.Context, vocabID uuid.UUID) (entity.Vocabulary, error) {
	query := `SELECT user_id, name, native_lang, translate_lang FROM vocabulary WHERE id=$1;`
	var vocab entity.Vocabulary
	err := r.pgxPool.QueryRow(ctx, query, vocabID).Scan(&vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByID: %w", err)
	}

	vocab.ID = vocabID

	return vocab, nil
}

func (r *VocabRepo) GetVocabularies(ctx context.Context, userID uuid.UUID) ([]entity.Vocabulary, error) {
	query := `SELECT v.id, v.user_id, name, n.lang as native_lang, t.lang as translate_lang, array_agg(tg."text") as tags, v.access FROM vocabulary v
LEFT JOIN "language" n ON n.code = v.native_lang
LEFT JOIN "language" t ON t.code = v.translate_lang 
LEFT JOIN "tag" tg ON tg.id = any(v.tags)
WHERE user_id=$1
GROUP BY v.id, n.lang, t.lang;`
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
			pq.Array(&sqlTags),
			&vocab.Access,
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

func (r *VocabRepo) GetCountVocabularies(ctx context.Context, userID uuid.UUID) (int, error) {
	var countVocabularies int

	query := `SELECT COUNT(id) FROM vocabulary WHERE user_id=$1;`

	err := r.pgxPool.QueryRow(ctx, query, userID).Scan(&countVocabularies)
	if err != nil {
		return -1, err
	}
	return countVocabularies, nil
}

func (r *VocabRepo) Edit(ctx context.Context, vocab entity.Vocabulary) error {
	query := `UPDATE vocabulary SET name=$2, access=$3 WHERE id=$1;`
	result, err := r.pgxPool.Exec(ctx, query, vocab.ID, vocab.Name, vocab.Access)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Edit: %w", err)
	}
	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Edit: %w", entity.ErrVocabularyNotFound)
	}
	return nil
}
