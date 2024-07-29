package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sorted "github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VocabRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *VocabRepo {
	return &VocabRepo{
		pgxPool: pgxPool,
	}
}

func (r *VocabRepo) AddVocab(ctx context.Context, vocab entity.Vocabulary, tagIDs []uuid.UUID) (uuid.UUID, error) {
	query := `
	INSERT INTO vocabulary (
		id, 
		user_id, 
		name, 
		native_lang, 
		translate_lang, 
		description, 
		tags, 
		updated_at, 
		created_at, 
		access) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $8, $9);`

	vid := uuid.New()
	_, err := r.pgxPool.Exec(ctx, query, vid, vocab.UserID, vocab.Name, vocab.NativeLang, vocab.TranslateLang, vocab.Description, tagIDs, time.Now().UTC(), vocab.Access)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Add: %w", err)
	}

	return vid, nil
}

func (r *VocabRepo) DeleteVocab(ctx context.Context, vocab entity.Vocabulary) error {
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

func (r *VocabRepo) GetVocab(ctx context.Context, vid uuid.UUID) (entity.Vocabulary, error) {
	query := `
	SELECT 
		id, 
		user_id, 
		name, 
		native_lang, 
		translate_lang, 
		description, 
		tags, 
		access, 
		created_at, 
		updated_at
	FROM vocabulary v
	WHERE id=$1;`
	var vocab entity.Vocabulary
	var tags []uuid.UUID
	err := r.pgxPool.QueryRow(ctx, query, vid).Scan(
		&vocab.ID,
		&vocab.UserID,
		&vocab.Name,
		&vocab.NativeLang,
		&vocab.TranslateLang,
		&vocab.Description,
		&tags,
		&vocab.Access,
		&vocab.CreatedAt,
		&vocab.UpdatedAt)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Get: %w", err)
	}

	vocab.Tags = make([]entityTag.Tag, 0, len(tags))
	for _, tagID := range tags {
		vocab.Tags = append(vocab.Tags, entityTag.Tag{ID: tagID})
	}

	return vocab, nil
}

func (r *VocabRepo) GetCreatorVocab(ctx context.Context, vocabID uuid.UUID) (uuid.UUID, error) {
	query := `
	SELECT user_id
	FROM vocabulary
	WHERE id=$1;`
	var userID uuid.UUID
	err := r.pgxPool.QueryRow(ctx, query, vocabID).Scan(&userID)
	if err != nil {
		return userID, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetCreatorVocab: %w", err)
	}
	return userID, nil
}

func (r *VocabRepo) GetByName(ctx context.Context, userID uuid.UUID, name string) (entity.Vocabulary, error) {
	query := `
	SELECT id, user_id, name, native_lang, translate_lang, description 
	FROM vocabulary v
	WHERE user_id=$1 AND name=$2;`
	var vocab entity.Vocabulary
	err := r.pgxPool.QueryRow(ctx, query, userID, name).Scan(&vocab.ID, &vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang, &vocab.Description)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetByName: %w", err)
	}
	return vocab, nil
}

func (r *VocabRepo) GetTagsVocabulary(ctx context.Context, vocabID uuid.UUID) ([]string, error) {
	query := `
	SELECT "text"
	FROM tag t
	LEFT JOIN vocabulary v ON t.id = ANY(v.tags)
	WHERE v.id=$1;`
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

func (r *VocabRepo) GetCountVocabularies(ctx context.Context, userID uuid.UUID) (int, error) {
	var countVocabularies int

	query := `SELECT COUNT(id) FROM vocabulary WHERE user_id=$1;`

	err := r.pgxPool.QueryRow(ctx, query, userID).Scan(&countVocabularies)
	if err != nil {
		return -1, err
	}
	return countVocabularies, nil
}

func (r *VocabRepo) EditVocab(ctx context.Context, vocab entity.Vocabulary) error {
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

func (r *VocabRepo) GetVocabulariesCountByAccess(ctx context.Context, uid uuid.UUID, accessIDs []uint8, search, nativeLang, translateLang string) (int, error) {
	query := fmt.Sprintf(`
	SELECT count(v.id)
	FROM vocabulary v
	WHERE (v.user_id=$1 OR v.access = ANY($2)) 
		AND (POSITION($3 in v."name")>0 OR POSITION($3 in v."description")>0) %s %s;`, getEqualLanguage("native_lang", nativeLang), getEqualLanguage("translate_lang", translateLang))

	var countLine int
	err := r.pgxPool.QueryRow(ctx, query, uid, accessIDs, search).Scan(&countLine)
	if err != nil {
		return 0, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetVocabulariesCountByAccess: %w", err)
	}

	return countLine, nil
}

func (r *VocabRepo) GetVocabulariesByAccess(ctx context.Context, uid uuid.UUID, accessIDs []uint8, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabularyWithUser, error) {
	query := fmt.Sprintf(`
	SELECT 
		v.id,
		v.user_id,
		u."name" "user_name",
		v.name,
		v.native_lang,
		v.translate_lang,
		v.description,
		count(w.id) as "words_count",
		array_agg(tg."text") as tags,
		v.access,
		v.updated_at,
		v.created_at
	FROM vocabulary v
	LEFT JOIN "tag" tg ON tg.id = any(v.tags)
	LEFT JOIN users u ON u.id = v.user_id
	LEFT JOIN word w ON w.vocabulary_id = v.id 
	WHERE (v.user_id=$1 OR v.access = ANY($2))
		AND (POSITION($3 in v."name")>0 OR POSITION($3 in v."description")>0) %s %s
	GROUP BY v.id, u."name"
	%s
	LIMIT $4
	OFFSET $5;`, getEqualLanguage("v.native_lang", nativeLang), getEqualLanguage("v.translate_lang", translateLang), getSorted(typeSort, sorted.TypeOrder(order)))
	rows, err := r.pgxPool.Query(ctx, query, uid, accessIDs, search, itemsPerPage, (page-1)*itemsPerPage)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetVocabulariesByAccess: %w", err)
	}
	defer rows.Close()

	var vocabularies []entity.VocabularyWithUser
	for rows.Next() {
		var vocab entity.VocabularyWithUser
		var sqlTags []sql.NullString
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.UserName,
			&vocab.Name,
			&vocab.NativeLang,
			&vocab.TranslateLang,
			&vocab.Description,
			&vocab.WordsCount,
			&sqlTags,
			&vocab.Access,
			&vocab.UpdatedAt,
			&vocab.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetVocabulariesByAccess - scan: %w", err)
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

func (r *VocabRepo) GetAccess(ctx context.Context, vid uuid.UUID) (uint8, error) {
	var accessID uint8
	err := r.pgxPool.QueryRow(ctx, "SELECT access FROM vocabulary WHERE id=$1", vid).Scan(&accessID)
	if err != nil {
		return 0, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.GetAccess: %w", err)
	}
	return accessID, nil
}

func (r *VocabRepo) CopyVocab(ctx context.Context, uid, vid uuid.UUID) (uuid.UUID, error) {
	vocab, err := r.GetVocab(ctx, vid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Copy - get vocab: %w", err)
	}

	tags := make([]uuid.UUID, 0, len(vocab.Tags))
	for _, t := range vocab.Tags {
		tags = append(tags, t.ID)
	}

	vocab.UserID = uid
	vid, err = r.AddVocab(ctx, vocab, tags)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Copy - add vocab: %w", err)
	}

	return vid, nil
}

func getSorted(typeSorted int, order sorted.TypeOrder) string {
	switch sorted.TypeSorted(typeSorted) {
	case sorted.Created:
		return fmt.Sprintf("ORDER BY v.created_at %s", order)
	case sorted.Updated:
		return fmt.Sprintf("ORDER BY v.updated_at %s", order)
	case sorted.ABC:
		return fmt.Sprintf("ORDER BY v.name %s", order)
	default:
		return ""
	}
}

func getEqualLanguage(field, lang string) string {
	switch lang {
	case "any":
		return ""
	default:
		return fmt.Sprintf("AND %s='%s'", field, lang)
	}
}
