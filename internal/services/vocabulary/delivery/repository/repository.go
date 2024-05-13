package repository

import (
	"context"
	"database/sql"
	"fmt"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	"time"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type VocabRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *VocabRepo {
	return &VocabRepo{
		db: db,
	}
}

func (r *VocabRepo) Add(ctx context.Context, vocab entity.Vocabulary, tagIDs []uuid.UUID) error {
	query := `INSERT INTO vocabulary (id, user_id, name, native_lang, translate_lang, tags, updated_at, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, $7)`

	_, err := r.db.ExecContext(ctx, query, vocab.ID, vocab.UserID, vocab.Name, vocab.NativeLang, vocab.TranslateLang, pq.Array(tagIDs), time.Now().UTC())
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Add: %w", err)
	}

	return nil
}

func (r *VocabRepo) Delete(ctx context.Context, vocab entity.Vocabulary) error {
	query := `DELETE FROM vocabulary WHERE user_id=$1 AND name=$2;`
	result, err := r.db.ExecContext(ctx, query, vocab.UserID, vocab.Name)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Delete: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Delete: %w", entity.ErrVocabularyNotFound)
	}

	return nil
}

func (r *VocabRepo) GetByName(ctx context.Context, userID uuid.UUID, name string) (entity.Vocabulary, error) {
	query := `SELECT id, user_id, name, n.lang as native_lang, t.lang as translate_lang FROM vocabulary v
left join "language" n on n.code =v.native_lang 
left join "language" t on t.code =v.translate_lang 
WHERE user_id=$1 AND name=$2;`
	var vocab entity.Vocabulary
	err := r.db.QueryRowContext(ctx, query, userID, name).Scan(&vocab.ID, &vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
	}
	return vocab, nil
}

func (r *VocabRepo) GetTagsVocabulary(ctx context.Context, vocabID uuid.UUID) ([]string, error) {
	query := `SELECT "text" from tag t
left join vocabulary v on t.id=any(v.tags)
where v.id=$1;`
	rows, err := r.db.QueryContext(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
	}
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
	err := r.db.QueryRowContext(ctx, query, vocabID).Scan(&vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByID: %w", err)
	}

	vocab.ID = vocabID

	return vocab, nil
}

func (r *VocabRepo) GetVocabularies(ctx context.Context, userID uuid.UUID) ([]entity.Vocabulary, error) {
	query := `SELECT v.id, v.user_id, name, n.lang as native_lang, t.lang as translate_lang, array_agg(tg."text") as tags FROM vocabulary v
LEFT JOIN "language" n ON n.code = v.native_lang
LEFT JOIN "language" t ON t.code = v.translate_lang 
LEFT JOIN "tag" tg ON tg.id = any(v.tags)
WHERE user_id=$1
GROUP BY v.id, n.lang, t.lang;`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

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

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&countVocabularies)
	if err != nil {
		return -1, err
	}
	return countVocabularies, nil
}

func (r *VocabRepo) Rename(ctx context.Context, id uuid.UUID, newName string) error {
	query := `UPDATE vocabulary SET name=$1 WHERE id=$2;`
	result, err := r.db.ExecContext(ctx, query, newName, id)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Rename: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Rename: %w", entity.ErrVocabularyNotFound)
	}
	return nil
}
