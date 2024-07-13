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

func (r *VocabRepo) Add(ctx context.Context, vocab entity.Vocabulary, tagIDs []uuid.UUID) error {
	query := `INSERT INTO vocabulary (id, user_id, name, native_lang, translate_lang, tags, updated_at, created_at, access, access_edit) VALUES($1, $2, $3, $4, $5, $6, $7, $7, $8, $9)`

	_, err := r.pgxPool.Exec(ctx, query, vocab.ID, vocab.UserID, vocab.Name, vocab.NativeLang, vocab.TranslateLang, tagIDs, time.Now().UTC(), vocab.Access, false)
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
	query := `
	SELECT id, user_id, name, native_lang, translate_lang 
	FROM vocabulary v
	WHERE id=$1;`
	var vocab entity.Vocabulary
	err := r.pgxPool.QueryRow(ctx, query, vocabID).Scan(&vocab.ID, &vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.delivery.repository.VocabRepo.Get: %w", err)
	}
	return vocab, nil
}

func (r *VocabRepo) GetByName(ctx context.Context, userID uuid.UUID, name string) (entity.Vocabulary, error) {
	query := `
	SELECT id, user_id, name, native_lang, translate_lang 
	FROM vocabulary v
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
	LEFT JOIN "tag" tg ON tg.id = any(v.tags)
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

func (r *VocabRepo) GetVocabulariesCountByAccess(ctx context.Context, uid uuid.UUID, accessIDs []int, search, nativeLang, translateLang string) (int, error) {
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

func (r *VocabRepo) GetVocabulariesByAccess(ctx context.Context, uid uuid.UUID, accessIDs []int, page, itemsPerPage, typeOrder int, search, nativeLang, translateLang string) ([]entity.VocabularyWithUser, error) {
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
	OFFSET $5;`, getEqualLanguage("v.native_lang", nativeLang), getEqualLanguage("v.translate_lang", translateLang), getSorted(typeOrder))
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

func getSorted(typeSorted int) string {
	switch sorted.TypeSorted(typeSorted) {
	case sorted.Newest:
		return "ORDER BY v.created_at DESC"
	case sorted.Oldest:
		return "ORDER BY v.created_at ASC"
	case sorted.UpdateAsc:
		return "ORDER BY v.updated_at ASC"
	case sorted.UpdateDesc:
		return "ORDER BY v.updated_at DESC"
	case sorted.AtoZ:
		return "ORDER BY v.name ASC"
	case sorted.ZtoA:
		return "ORDER BY v.name DESC"
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
