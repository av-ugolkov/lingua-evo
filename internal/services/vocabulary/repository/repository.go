package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	sorted "github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
)

type VocabRepo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *VocabRepo {
	return &VocabRepo{
		tr: tr,
	}
}

func (r *VocabRepo) AddVocab(ctx context.Context, vocab entity.Vocab) (uuid.UUID, error) {
	query := `
	INSERT INTO vocabulary (
		id, 
		user_id, 
		name, 
		native_lang, 
		translate_lang, 
		description, 
		updated_at, 
		created_at, 
		access) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $7, $8);`

	vid := uuid.New()
	_, err := r.tr.Exec(ctx, query, vid, vocab.UserID, vocab.Name, vocab.NativeLang, vocab.TranslateLang, vocab.Description, time.Now().UTC(), vocab.Access)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.repository.VocabRepo.Add: %w", err)
	}

	return vid, nil
}

func (r *VocabRepo) DeleteVocab(ctx context.Context, vocab entity.Vocab) error {
	query := `DELETE FROM vocabulary WHERE user_id=$1 AND name=$2;`
	result, err := r.tr.Exec(ctx, query, vocab.UserID, vocab.Name)
	if err != nil {
		return fmt.Errorf("vocabulary.repository.VocabRepo.Delete: %w", err)
	}
	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("vocabulary.repository.VocabRepo.Delete: %w", entity.ErrVocabularyNotFound)
	}

	return nil
}

func (r *VocabRepo) GetVocab(ctx context.Context, vid uuid.UUID) (entity.Vocab, error) {
	query := `
	SELECT 
		id, 
		user_id, 
		name, 
		native_lang, 
		translate_lang, 
		description, 
		access, 
		created_at, 
		updated_at
	FROM vocabulary v
	WHERE id=$1;`
	var vocab entity.Vocab
	err := r.tr.QueryRow(ctx, query, vid).Scan(
		&vocab.ID,
		&vocab.UserID,
		&vocab.Name,
		&vocab.NativeLang,
		&vocab.TranslateLang,
		&vocab.Description,
		&vocab.Access,
		&vocab.CreatedAt,
		&vocab.UpdatedAt)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.repository.VocabRepo.Get: %w", err)
	}

	return vocab, nil
}

func (r *VocabRepo) GetCreatorVocab(ctx context.Context, vocabID uuid.UUID) (uuid.UUID, error) {
	query := `
	SELECT user_id
	FROM vocabulary
	WHERE id=$1;`
	var userID uuid.UUID
	err := r.tr.QueryRow(ctx, query, vocabID).Scan(&userID)
	if err != nil {
		return userID, fmt.Errorf("vocabulary.repository.VocabRepo.GetCreatorVocab: %w", err)
	}
	return userID, nil
}

func (r *VocabRepo) GetByName(ctx context.Context, userID uuid.UUID, name string) (entity.Vocab, error) {
	query := `
	SELECT id, user_id, name, native_lang, translate_lang, description 
	FROM vocabulary v
	WHERE user_id=$1 AND name=$2;`
	var vocab entity.Vocab
	err := r.tr.QueryRow(ctx, query, userID, name).Scan(&vocab.ID, &vocab.UserID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang, &vocab.Description)
	if err != nil {
		return vocab, fmt.Errorf("vocabulary.repository.VocabRepo.GetByName: %w", err)
	}
	return vocab, nil
}

func (r *VocabRepo) GetTagsVocabulary(ctx context.Context, vocabID uuid.UUID) ([]string, error) {
	query := `
	SELECT "text"
	FROM tag t
	LEFT JOIN vocabulary v ON t.id = ANY(v.tags)
	WHERE v.id=$1;`
	rows, err := r.tr.Query(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetByName: %w", err)
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetByName: %w", err)
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *VocabRepo) GetCountVocabularies(ctx context.Context, userID uuid.UUID) (int, error) {
	var countVocabularies int

	query := `SELECT COUNT(id) FROM vocabulary WHERE user_id=$1;`

	err := r.tr.QueryRow(ctx, query, userID).Scan(&countVocabularies)
	if err != nil {
		return -1, err
	}
	return countVocabularies, nil
}

func (r *VocabRepo) EditVocab(ctx context.Context, vocab entity.Vocab) error {
	query := `UPDATE vocabulary SET name=$2, access=$3, description=$4 WHERE id=$1;`
	result, err := r.tr.Exec(ctx, query, vocab.ID, vocab.Name, vocab.Access, vocab.Description)
	if err != nil {
		return fmt.Errorf("vocabulary.repository.VocabRepo.Edit: %w", err)
	}
	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("vocabulary.repository.VocabRepo.Edit: %w", entity.ErrVocabularyNotFound)
	}
	return nil
}

func (r *VocabRepo) GetVocabulariesCountByAccess(ctx context.Context, uid uuid.UUID, accessTypes []access.Type, search, nativeLang, translateLang string) (int, error) {
	query := fmt.Sprintf(`
	SELECT count(v.id)
	FROM vocabulary v
	WHERE (v.user_id=$1 OR v.access = ANY($2)) 
		AND (v."name" LIKE '%[1]s' || $3 || '%[1]s' OR v."description" LIKE '%[1]s' || $3 || '%[1]s') %[2]s %[3]s;`, "%",
		getEqualLanguage("native_lang", nativeLang),
		getEqualLanguage("translate_lang", translateLang))

	var countLine int
	err := r.tr.QueryRow(ctx, query, uid, accessTypes, search).Scan(&countLine)
	if err != nil {
		return 0, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesCountByAccess: %w", err)
	}

	return countLine, nil
}

func (r *VocabRepo) GetVocabulariesByAccess(ctx context.Context, uid uuid.UUID, accessTypes []access.Type, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabWithUser, error) {
	query := fmt.Sprintf(`
	SELECT 
		v.id,
		v.user_id,
		u."nickname",
		v.name,
		v.native_lang,
		v.translate_lang,
		v.description,
		count(w.id) as "words_count",
		v.access,
		v.updated_at,
		v.created_at
	FROM vocabulary v
	LEFT JOIN users u ON u.id = v.user_id
	LEFT JOIN word w ON w.vocabulary_id = v.id 
	WHERE (v.user_id=$1 OR v.access = ANY($2))
		AND (v."name" LIKE '%[1]s' || $3 || '%[1]s' OR v."description" LIKE '%[1]s' || $3 || '%[1]s') %[2]s %[3]s
	GROUP BY v.id, u."nickname"
	%[4]s
	LIMIT $4
	OFFSET $5;`, "%", getEqualLanguage("v.native_lang", nativeLang), getEqualLanguage("v.translate_lang", translateLang), getSorted(typeSort, sorted.TypeOrder(order)))
	rows, err := r.tr.Query(ctx, query, uid, accessTypes, search, itemsPerPage, (page-1)*itemsPerPage)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesByAccess: %w", err)
	}
	defer rows.Close()

	vocabularies := make([]entity.VocabWithUser, 0, itemsPerPage)
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
			return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesByAccess - scan: %w", err)
		}

		vocabularies = append(vocabularies, vocab)
	}

	return vocabularies, nil
}

func (r *VocabRepo) GetAccess(ctx context.Context, vid uuid.UUID) (uint8, error) {
	var accessID uint8
	err := r.tr.QueryRow(ctx, "SELECT access FROM vocabulary WHERE id=$1", vid).Scan(&accessID)
	if err != nil {
		return 0, fmt.Errorf("vocabulary.repository.VocabRepo.GetAccess: %w", err)
	}
	return accessID, nil
}

func (r *VocabRepo) CopyVocab(ctx context.Context, uid, vid uuid.UUID) (uuid.UUID, error) {
	vocab, err := r.GetVocab(ctx, vid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.repository.VocabRepo.Copy - get vocab: %w", err)
	}

	vocab.UserID = uid
	vid, err = r.AddVocab(ctx, vocab)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.repository.VocabRepo.Copy - add vocab: %w", err)
	}

	return vid, nil
}

func (r *VocabRepo) GetVocabsWithCountWords(ctx context.Context, uid, owner uuid.UUID, access []uint8) ([]entity.VocabWithUser, error) {
	var limit int
	err := r.tr.QueryRow(ctx, `
		SELECT count(v.id) FROM vocabulary v
		WHERE v.user_id = $1 AND "access" = any($2)`, owner, access).Scan(&limit)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabsWithCountWords - get limit: %w", err)
	}

	query := `
		SELECT 
		    v.id, 
		    name, 
		    native_lang, 
		    translate_lang, 
		    access, 
		    count(w.id), 
		    count(vn.user_id)!=0 notification 
		FROM vocabulary v
		LEFT JOIN word w ON w.vocabulary_id = v.id
		LEFT JOIN vocabulary_notifications vn ON vn.user_id=$3 AND vn.vocab_id=v.id
		WHERE v.user_id = $1 AND "access" = any($2)
		GROUP BY v.id
		LIMIT $4;`

	rows, err := r.tr.Query(ctx, query, owner, access, uid, limit)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetWithCountWords: %w", err)
	}

	vocabs := make([]entity.VocabWithUser, 0, limit)
	var vocab entity.VocabWithUser
	for rows.Next() {
		err := rows.Scan(&vocab.ID, &vocab.Name, &vocab.NativeLang, &vocab.TranslateLang, &vocab.Access, &vocab.WordsCount, &vocab.Notification)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetWithCountWords - scan: %w", err)
		}
		vocabs = append(vocabs, vocab)
	}

	return vocabs, nil
}

func (r *VocabRepo) GetWithCountWords(ctx context.Context, vid uuid.UUID) (entity.VocabWithUser, error) {
	query := `
		SELECT 
			v.id, 
			v.name,
			v.user_id,
			native_lang, 
			translate_lang, 
			access, 
			count(w.id), 
			v.description, 
			v.created_at, 
			v.updated_at, 
			u."nickname" 
		FROM vocabulary v
		LEFT JOIN word w ON w.vocabulary_id = v.id
		LEFT JOIN users u ON u.id = v.user_id 
		WHERE v.id = $1
		GROUP BY v.id, u."nickname";`

	var vocab entity.VocabWithUser
	err := r.tr.QueryRow(ctx, query, vid).Scan(
		&vocab.ID,
		&vocab.Name,
		&vocab.UserID,
		&vocab.NativeLang,
		&vocab.TranslateLang,
		&vocab.Access,
		&vocab.WordsCount,
		&vocab.Description,
		&vocab.CreatedAt,
		&vocab.UpdatedAt,
		&vocab.UserName,
	)
	if err != nil {
		return entity.VocabWithUser{}, fmt.Errorf("vocabulary.repository.VocabRepo.GetWithCountWords: %w", err)
	}

	return vocab, nil
}

func (r *VocabRepo) GetVocabulariesWithMaxWords(ctx context.Context, access []uint8, limit int) ([]entity.VocabWithUser, error) {
	const query = `
		SELECT 
			v.id,
			user_id, 
			name, 
			native_lang, 
			translate_lang, 
			access, 
			count(w.id) cw, 
			v.description 
		FROM vocabulary v
		LEFT JOIN word w ON w.vocabulary_id = v.id 
		WHERE v.access = ANY($2)
		GROUP BY v.id
		HAVING count(w.id) > 0
		ORDER BY cw DESC
		LIMIT $1`

	rows, err := r.tr.Query(ctx, query, limit, access)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesRandom: %w", err)
	}

	vocabs := make([]entity.VocabWithUser, 0, limit)
	var vocab entity.VocabWithUser
	for rows.Next() {
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.Name,
			&vocab.NativeLang,
			&vocab.TranslateLang,
			&vocab.Access,
			&vocab.WordsCount,
			&vocab.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesRandom - scan: %w", err)
		}
		vocabs = append(vocabs, vocab)
	}

	return vocabs, nil
}

func (r *VocabRepo) GetVocabulariesRecommended(ctx context.Context, uid uuid.UUID, access []uint8, limit uint) ([]entity.VocabWithUser, error) {
	const query = `
		SELECT 
			v.id,
			user_id,
			name,
			native_lang,
			translate_lang,
			access,
			count(w.id) cw,
			v.description
		FROM vocabulary v
		LEFT JOIN word w ON w.vocabulary_id = v.id 
		WHERE v.access = ANY($2) 
			AND v.native_lang = any(SELECT DISTINCT native_lang FROM vocabulary v WHERE v.user_id = $1) 
			AND v.user_id != $1
		GROUP BY v.id
		HAVING count(w.id) > 0
		ORDER BY cw DESC 
		LIMIT $3;`

	rows, err := r.tr.Query(ctx, query, uid, access, limit)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesRecommended: %w", err)
	}
	defer rows.Close()

	vocabs := make([]entity.VocabWithUser, 0, limit)
	var vocab entity.VocabWithUser
	for rows.Next() {
		err := rows.Scan(
			&vocab.ID,
			&vocab.UserID,
			&vocab.Name,
			&vocab.NativeLang,
			&vocab.TranslateLang,
			&vocab.Access,
			&vocab.WordsCount,
			&vocab.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.repository.VocabRepo.GetVocabulariesRecommended - scan: %w", err)
		}
		vocabs = append(vocabs, vocab)
	}

	return vocabs, nil
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
		return runtime.EmptyString
	}
}

func getEqualLanguage(field, lang string) string {
	switch lang {
	case "any":
		return runtime.EmptyString
	default:
		return fmt.Sprintf("AND %s='%s'", field, lang)
	}
}

func getDictTable(langCode string) string {
	table := "dictionary"
	if len(langCode) != 0 {
		table = fmt.Sprintf(`%s_%s`, table, langCode)
	}
	return table
}

func getExamTable(langCode string) string {
	table := "example"
	if len(langCode) != 0 {
		table = fmt.Sprintf(`%s_%s`, table, langCode)
	}
	return table
}
