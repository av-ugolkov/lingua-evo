package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

type VocabRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *VocabRepo {
	return &VocabRepo{
		db: db,
	}
}

func (r *VocabRepo) Add(ctx context.Context, vocab entity.Vocabulary) error {
	query := `INSERT INTO vocabulary (id, user_id, name, native_lang, translate_lang, updated_at, created_at) VALUES($1, $2, $3, $4, $5, $6, $6)`

	_, err := r.db.ExecContext(ctx, query, vocab.ID, vocab.UserID, vocab.Name, vocab.NativeLang, vocab.TranslateLang, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Add: %w", err)
	}

	return nil
}

func (r *VocabRepo) AddTagsToVocabulary(ctx context.Context, vocabularyID uuid.UUID, tagIDs []uuid.UUID) error {
	query := `INSERT INTO vocabulary_tag (vocabulary_id, tag_id) VALUES`
	vals := make([]interface{}, 0, len(tagIDs))
	vals = append(vals, vocabularyID)
	for ind, tagID := range tagIDs {
		query += fmt.Sprintf("($1, $%d),", ind+2)
		vals = append(vals, tagID)
	}
	query = query[0 : len(query)-1]
	_, err := r.db.ExecContext(ctx, query, vals...)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.AddTagsToVocabulary: %w", err)
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
	query := `SELECT id, user_id, name, n.lang native_lang, t.lang translate_lang FROM vocabulary v
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
	query := `SELECT "text" from tag 
where id=any(select tag_id from vocabulary_tag where vocabulary_id=$1);`
	rows, err := r.db.QueryContext(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
	}
	tags := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *VocabRepo) GetByID(ctx context.Context, dictID uuid.UUID) (entity.Vocabulary, error) {
	query := `SELECT user_id, name, native_lang, translate_lang FROM vocabulary WHERE id=$1;`
	var vocab entity.Vocabulary
	err := r.db.QueryRowContext(ctx, query, dictID).Scan(&vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByID: %w", err)
	}
	return vocab, nil
}

func (r *VocabRepo) GetVocabularies(ctx context.Context, userID uuid.UUID) ([]entity.Vocabulary, error) {
	query := `SELECT d.id, d.user_id, name, n.lang native_lang, s.lang translate_lang FROM vocabulary d
left join "language" n on n.code = d.native_lang
left join "language" s on s.code = d.translate_lang 
WHERE user_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vocabularies []entity.Vocabulary
	for rows.Next() {
		var vocab entity.Vocabulary
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.Name,
			&vocab.NativeLang,
			&vocab.TranslateLang,
		)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetVocabularies - scan: %w", err)
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
