package repository

import (
	"context"
	"fmt"

	sorted "github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
)

func (r *VocabRepo) GetVocabulariesByUser(ctx context.Context, userID uuid.UUID) ([]entity.VocabWithUser, error) {
	query := `
	SELECT v.id,
		v.user_id,
		v.name,
		v.native_lang,
		v.translate_lang,
		v.description,
		v.access,
		v.created_at,
		u.nickname,
		count(w.id) cw
	FROM vocabulary v
	LEFT JOIN users u ON u.id = v.user_id 
	LEFT JOIN word w ON w.vocabulary_id = v.id 
	WHERE user_id=$1
	GROUP BY v.id, u.nickname;`
	rows, err := r.tr.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesByUser: %w", err)
	}
	defer rows.Close()

	var vocabularies []entity.VocabWithUser
	var vocab entity.VocabWithUser
	for rows.Next() {
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.Name,
			&vocab.NativeLang,
			&vocab.TranslateLang,
			&vocab.Description,
			&vocab.Access,
			&vocab.CreatedAt,
			&vocab.UserName,
			&vocab.WordsCount)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabularies - scan: %w", err)
		}

		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}

func (r *VocabRepo) GetVocabulariesCountByUser(ctx context.Context, uid uuid.UUID, access []access.Type, search, nativeLang, translateLang string) (int, error) {
	query := fmt.Sprintf(`
	SELECT (count(v.id)+(select count(v.id) FROM vocabulary v 
	LEFT JOIN subscribers s ON s.user_id = $1
	WHERE v.user_id =s.subscribers_id AND v."access" = ANY($2))) AS count FROM vocabulary v
	WHERE v.user_id=$1
		AND (v."name" LIKE '%[1]s' || $3 || '%[1]s' OR v."description" LIKE '%[1]s' || $3 || '%[1]s') %[2]s %[3]s;`,
		"%",
		getEqualLanguage("native_lang", nativeLang),
		getEqualLanguage("translate_lang", translateLang))

	var countLine int
	err := r.tr.QueryRow(ctx, query, uid, access, search).Scan(&countLine)
	if err != nil {
		return 0, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesCountByAccess: %w", err)
	}

	return countLine, nil
}

func (r *VocabRepo) GetSortedVocabulariesByUser(ctx context.Context, userID uuid.UUID, accessTypes []access.Type, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabWithUser, error) {
	query := fmt.Sprintf(`
	WITH vocabulary_data AS (
    SELECT 
        v.id,
        v.user_id,
        u."nickname",
        v.name,
        v.native_lang,
        v.translate_lang,
        v.description,
        count(w.id) AS words_count,
        v.access,
        v.updated_at,
        v.created_at
    FROM vocabulary v
    LEFT JOIN users u ON u.id = v.user_id
    LEFT JOIN word w ON w.vocabulary_id = v.id
    LEFT JOIN subscribers s ON s.subscribers_id = v.user_id 
        AND s.user_id = $1 
    WHERE v.user_id = $1 
       OR (s.user_id IS NOT NULL AND v.access = ANY($2))
    GROUP BY v.id, u."nickname")
	SELECT *
	FROM vocabulary_data v
	WHERE (v."name" LIKE '%[1]s' || $3 || '%[1]s' OR v."description" LIKE '%[1]s' || $3 || '%[1]s') %[2]s %[3]s
	%[4]s
	LIMIT $4
	OFFSET $5;`,
		"%",
		getEqualLanguage("v.native_lang", nativeLang),
		getEqualLanguage("v.translate_lang", translateLang),
		getSorted(typeSort, sorted.TypeOrder(order)))
	rows, err := r.tr.Query(ctx, query, userID, accessTypes, search, itemsPerPage, (page-1)*itemsPerPage)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetSortedVocabulariesByUser: %w", err)
	}
	defer rows.Close()

	var vocabularies []entity.VocabWithUser
	var vocab entity.VocabWithUser
	for rows.Next() {
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.UserName,
			&vocab.Name,
			&vocab.NativeLang,
			&vocab.TranslateLang,
			&vocab.Description,
			&vocab.WordsCount,
			&vocab.Access,
			&vocab.UpdatedAt,
			&vocab.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetSortedVocabulariesByUser - scan: %w", err)
		}

		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}
