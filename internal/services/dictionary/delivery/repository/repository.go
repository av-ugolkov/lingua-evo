package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
)

type DictionaryRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *DictionaryRepo {
	return &DictionaryRepo{
		db: db,
	}
}

func (r *DictionaryRepo) AddWords(ctx context.Context, words []entity.DictWord) ([]uuid.UUID, error) {
	wordTexts := make([]string, 0, len(words))
	statements := make([]string, 0, len(words))
	params := make([]any, 0, len(words)+1)
	params = append(params, &wordTexts)
	counter := len(params)
	for _, word := range words {
		wordTexts = append(wordTexts, word.Text)
		statement := "$" + strconv.Itoa(counter+1) +
			",$" + strconv.Itoa(counter+2) +
			",$" + strconv.Itoa(counter+3) +
			",$" + strconv.Itoa(counter+4) +
			",$" + strconv.Itoa(counter+5) +
			",$" + strconv.Itoa(counter+6)

		counter += 6
		statements = append(statements, "("+statement+")")

		params = append(params, word.ID, word.Text, word.Pronunciation, word.LangCode,
			time.Now().UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))
	}

	table := getTable(words[0].LangCode)
	query := fmt.Sprintf(
		`WITH d AS (
    		SELECT id FROM %[1]s WHERE text = ANY($1::text[])),
		ins AS (
    		INSERT INTO %[1]s (id, text, pronunciation, lang_code, updated_at, created_at)
			VALUES %[2]s
    		ON CONFLICT DO NOTHING RETURNING id)
		SELECT id 
		FROM ins 
		UNION ALL 
		SELECT id 
		FROM d;`, table, strings.Join(statements, ", "))
	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.AddWord - query: %w", err)
	}

	wordIDs := make([]uuid.UUID, 0, len(words))
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.AddWord - scan: %w", err)
		}
		wordIDs = append(wordIDs, id)
	}

	return wordIDs, nil
}

func (r *DictionaryRepo) GetWordIDByText(ctx context.Context, w *entity.DictWord) (uuid.UUID, error) {
	table := getTable(w.LangCode)
	query := fmt.Sprintf(`SELECT id FROM "%s" WHERE text=$1 AND lang_code=$2;`, table)
	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, query, w.Text, w.LangCode).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWordByText: %w", err)
	}
	return id, nil
}

func (r *DictionaryRepo) GetWords(ctx context.Context, ids []uuid.UUID) ([]entity.DictWord, error) {
	query := `SELECT id, text, pronunciation FROM dictionary WHERE id=ANY($1);`
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWords - query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	words := make([]entity.DictWord, 0, len(ids))
	for rows.Next() {
		var word entity.DictWord
		err = rows.Scan(&word.ID, &word.Text, &word.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWords - scan: %w", err)
		}
		words = append(words, word)
	}
	return words, nil
}

func (r *DictionaryRepo) UpdateWord(ctx context.Context, w *entity.DictWord) error {
	query := `UPDATE dictionary SET text=$1, pronunciation=$2 WHERE id=$3`
	result, err := r.db.ExecContext(ctx, query, w.Text, w.Pronunciation, w.ID)
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictionaryRepo.EditWord - exec: %w", err)
	}

	if rows, err := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("dictionary.repository.DictionaryRepo.EditWord: not fount effected rows")
	} else if err != nil {
		return fmt.Errorf("dictionary.repository.DictionaryRepo.EditWord - rows affected: %w", err)
	}

	return nil
}

func (r *DictionaryRepo) FindWords(ctx context.Context, w *entity.DictWord) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	query := `SELECT id FROM dictionary WHERE text=$1 AND lang_code=$2;`
	rows, err := r.db.QueryContext(ctx, query, w.Text, w.LangCode)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.FindWords - query: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

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

func (r *DictionaryRepo) DeleteWordByID(ctx context.Context, id uuid.UUID) (int64, error) {
	query := `DELETE FROM dictionary WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWordByID - exec: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWordByID - rows affected: %w", err)
	}
	return rows, nil
}

func (r *DictionaryRepo) DeleteWordByText(ctx context.Context, text, langCode string) (int64, error) {
	query := `DELETE FROM dictionary WHERE text=$1 AND lang_code=$2`
	result, err := r.db.ExecContext(ctx, query, text, langCode)
	if err != nil {
		return 0, fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWordByText - exec: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWordByText - rows affected: %w", err)
	}
	return rows, nil
}

func (r *DictionaryRepo) GetRandomWord(ctx context.Context, langCode string) (entity.DictWord, error) {
	table := getTable(langCode)
	query := fmt.Sprintf(`SELECT text, pronunciation, lang_code FROM "%s" WHERE moderator IS NOT NULL ORDER BY RANDOM() LIMIT 1;`, table)
	word := entity.DictWord{}
	err := r.db.QueryRowContext(ctx, query).Scan(&word.Text, &word.Pronunciation, &word.LangCode)
	if err != nil {
		return entity.DictWord{}, fmt.Errorf("dictionary.repository.DictionaryRepo.GetRandomWord - scan: %w", err)
	}

	return word, nil
}

func getTable(langCode string) string {
	table := "dictionary"
	if len(langCode) != 0 {
		table = fmt.Sprintf(`%s_%s`, table, langCode)
	}
	return table
}
