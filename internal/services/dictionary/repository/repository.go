package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"

	"github.com/google/uuid"
)

type DictionaryRepo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *DictionaryRepo {
	return &DictionaryRepo{
		tr: tr,
	}
}

func (r *DictionaryRepo) GetCountDictionaryWords(ctx context.Context, langCode string) (int, error) {
	query := fmt.Sprintf(`
		SELECT 
			count(id) 
		FROM dictionary_%[1]s;`, langCode)

	var count int
	err := r.tr.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("dictionary.repository.DictionaryRepo.GetDictionary: %w", err)
	}
	return count, nil
}

func (r *DictionaryRepo) GetDictionary(ctx context.Context, langCode, search string, page, itemsPerPage int) ([]entity.DictWordData, error) {
	query := fmt.Sprintf(`
		SELECT 
			d.id, 
			text, 
			COALESCE(pronunciation, ''), 
			COALESCE(u.nickname, ''),
			lang_code,
			d.created_at
		FROM dictionary_%[1]s d
		LEFT JOIN users u ON d.creator = u.id 
		WHERE text ILIKE '%[2]s' || $1 || '%[2]s'
		ORDER BY text
		LIMIT $2
		OFFSET $3;`, langCode, "%")

	rows, err := r.tr.Query(ctx, query, search, itemsPerPage, (page-1)*itemsPerPage)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetDictionary: %w", err)
	}
	defer rows.Close()

	words := make([]entity.DictWordData, 0, 100)
	var word entity.DictWordData
	for rows.Next() {
		err = rows.Scan(
			&word.ID,
			&word.Text,
			&word.Pronunciation,
			&word.Creator,
			&word.LangCode,
			&word.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetDictionary: %w", err)
		}

		words = append(words, word)
	}

	return words, nil
}

func (r *DictionaryRepo) AddWords(ctx context.Context, inWords []entity.DictWord) ([]entity.DictWord, error) {
	statements := make([]string, 0, len(inWords))
	params := make([]any, 0, len(inWords))
	counter := len(params)
	now := time.Now().UTC()
	for _, word := range inWords {
		statement := "$" + strconv.Itoa(counter+1) +
			",$" + strconv.Itoa(counter+2) +
			",$" + strconv.Itoa(counter+3) +
			",$" + strconv.Itoa(counter+4) +
			",$" + strconv.Itoa(counter+5) +
			",$" + strconv.Itoa(counter+6)

		counter += 6
		statements = append(statements, "("+statement+")")

		params = append(params, uuid.New(), word.Text, word.LangCode, word.Creator, now.Format(time.RFC3339), now.Format(time.RFC3339))
	}

	table := getTable(inWords[0].LangCode)
	query := fmt.Sprintf(`
		INSERT INTO "%s" (
			id, 
			text, 
			lang_code, 
			creator, 
			updated_at, 
			created_at) 
		VALUES %s ON CONFLICT DO NOTHING RETURNING 
			id, 
			text, 
			lang_code, 
			creator, 
			updated_at, 
			created_at;`, table, strings.Join(statements, ", "))
	rows, err := r.tr.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.AddWords - query: %w", err)
	}
	defer rows.Close()

	words := make([]entity.DictWord, 0, len(inWords))
	for rows.Next() {
		var word entity.DictWord
		if err := rows.Scan(&word.ID, &word.Text, &word.LangCode, &word.Creator, &word.UpdatedAt, &word.CreatedAt); err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.AddWords - scan: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

func (r *DictionaryRepo) GetWordsByText(ctx context.Context, inWords []entity.DictWord) ([]entity.DictWord, error) {
	texts := make([]string, 0, len(inWords))
	for _, word := range inWords {
		texts = append(texts, word.Text)
	}
	table := getTable(inWords[0].LangCode)
	query := fmt.Sprintf(`
		SELECT 
			id, 
			text, 
			pronunciation, 
			lang_code, 
			creator, 
			updated_at, 
			created_at 
		FROM "%s" 
		WHERE text = ANY($1::text[]);`, table)

	rows, err := r.tr.Query(ctx, query, texts)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWordsByText: %w", err)
	}
	defer rows.Close()

	words := make([]entity.DictWord, 0, len(inWords))
	var pron *string
	for rows.Next() {
		var word entity.DictWord
		if err := rows.Scan(
			&word.ID,
			&word.Text,
			&pron,
			&word.LangCode,
			&word.Creator,
			&word.UpdatedAt,
			&word.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWordsByText: %w", err)
		}
		if pron != nil {
			word.Pronunciation = *pron
		}
		words = append(words, word)
	}

	return words, nil
}

func (r *DictionaryRepo) GetWords(ctx context.Context, ids []uuid.UUID) ([]entity.DictWord, error) {
	query := `
		SELECT 
			id, 
			text, 
			pronunciation, 
			lang_code, 
			creator, 
			updated_at, 
			created_at 
		FROM dictionary 
		WHERE id=ANY($1);`
	rows, err := r.tr.Query(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWords - query: %w", err)
	}
	defer rows.Close()

	words := make([]entity.DictWord, 0, len(ids))
	var pron *string
	for rows.Next() {
		var word entity.DictWord
		err = rows.Scan(
			&word.ID,
			&word.Text,
			&pron,
			&word.LangCode,
			&word.Creator,
			&word.UpdatedAt,
			&word.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWords - scan: %w", err)
		}

		if pron != nil {
			word.Pronunciation = *pron
		}

		words = append(words, word)
	}
	return words, nil
}

func (r *DictionaryRepo) UpdateWord(ctx context.Context, w *entity.DictWord) error {
	query := `UPDATE dictionary SET text=$2, pronunciation=$3, updated_at=$4 WHERE id=$1`
	result, err := r.tr.Exec(ctx, query, w.ID, w.Text, w.Pronunciation, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictionaryRepo.EditWord - exec: %w", err)
	}

	if rows := result.RowsAffected(); rows > 1 {
		return fmt.Errorf("dictionary.repository.DictionaryRepo.EditWord: %w", entity.ErrAffectRows)
	}

	return nil
}

func (r *DictionaryRepo) FindWords(ctx context.Context, w *entity.DictWord) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	query := `SELECT id FROM dictionary WHERE text=$1 AND lang_code=$2;`
	rows, err := r.tr.Query(ctx, query, w.Text, w.LangCode)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.FindWords - query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.FindWords - scan: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (r *DictionaryRepo) DeleteWordByText(ctx context.Context, w *entity.DictWord) error {
	query := `DELETE FROM dictionary WHERE text=$1 AND lang_code=$2`
	result, err := r.tr.Exec(ctx, query, w.Text, w.LangCode)
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWordByText - exec: %w", err)
	}

	if rows := result.RowsAffected(); rows > 1 {
		return fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWordByText - more than 1 rows deleted")
	}
	return nil
}

func (r *DictionaryRepo) GetRandomWords(ctx context.Context, langCode string, count int) ([]entity.DictWord, error) {
	const query = `
		SELECT 
			id, 
			text, 
			COALESCE(pronunciation, ''), 
			lang_code 
		FROM "dictionary" 
		ORDER BY RANDOM() 
		LIMIT $1;`
	rows, err := r.tr.Query(ctx, query, count)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetRandomWord: %w", err)
	}
	defer rows.Close()

	var word entity.DictWord
	words := make([]entity.DictWord, 0, count)
	for rows.Next() {
		err = rows.Scan(
			&word.ID,
			&word.Text,
			&word.Pronunciation,
			&word.LangCode,
		)
		if err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetRandomWord: %w", err)
		}

		words = append(words, word)
	}

	return words, nil
}

func (r *DictionaryRepo) GetPronunciation(ctx context.Context, text, langCode string) (string, error) {
	table := getTable(langCode)

	query := fmt.Sprintf(`
		SELECT COALESCE(pronunciation, '') 
		FROM "%s" 
		WHERE text=$1;`, table)

	var pron string
	err := r.tr.QueryRow(ctx, query, text).Scan(&pron)
	if err != nil {
		return "", fmt.Errorf("dictionary.repository.DictionaryRepo.GetPronunciation: %w", err)
	}
	return pron, nil
}

func getTable(langCode string) string {
	table := "dictionary"
	if len(langCode) != 0 {
		table = fmt.Sprintf(`%s_%s`, table, langCode)
	}
	return table
}
