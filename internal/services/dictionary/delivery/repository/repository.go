package repository

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *DictionaryRepo) AddWords(ctx context.Context, words []entity.Word) ([]uuid.UUID, error) {
	if len(words) == 0 {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.AddWord - empty words")
	}

	wordTexts := make([]string, 0, len(words))
	var insValues strings.Builder
	for i := 0; i < len(words); i++ {
		wordTexts = append(wordTexts, words[i].Text)
		insValues.WriteString(fmt.Sprintf("('%v','%s','%s','%s','%q','%q')",
			words[i].ID,
			words[i].Text,
			words[i].Pronunciation,
			words[i].LangCode,
			time.Now().UTC().Format(time.RFC3339),
			time.Now().UTC().Format(time.RFC3339),
		))
		if i < len(words)-1 {
			insValues.WriteString(",")
		}
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
		FROM d;`, table, insValues.String())
	rows, err := r.db.QueryContext(ctx, query, wordTexts)
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

func (r *DictionaryRepo) GetWordByText(ctx context.Context, w *entity.Word) (uuid.UUID, error) {
	word := &entity.Word{}
	table := getTable(w.LangCode)
	query := fmt.Sprintf(`SELECT id FROM "%s" WHERE text=$1 AND lang_code=$2;`, table)
	err := r.db.QueryRowContext(ctx, query, w.Text, w.LangCode).Scan(&word.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWordByText: %w", err)
	}
	return word.ID, nil
}

func (r *DictionaryRepo) GetWords(ctx context.Context, ids []uuid.UUID) ([]entity.Word, error) {
	query := `SELECT id, text, pronunciation FROM dictionary WHERE id=ANY($1);`
	rows, err := r.db.QueryContext(ctx, query, ids)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWords - query: %w", err)
	}
	defer rows.Close()

	words := make([]entity.Word, 0, len(ids))
	for rows.Next() {
		var word entity.Word
		err = rows.Scan(&word.ID, &word.Text, &word.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetWords - scan: %w", err)
		}
		words = append(words, word)
	}
	return words, nil
}

func (r *DictionaryRepo) UpdateWord(ctx context.Context, w *entity.Word) error {
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

func (r *DictionaryRepo) FindWords(ctx context.Context, w *entity.Word) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	query := `SELECT id FROM dictionary WHERE text=$1% AND lang_code=$2;`
	rows, err := r.db.QueryContext(ctx, query, w.Text, w.LangCode)
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

func (r *DictionaryRepo) DeleteWord(ctx context.Context, w *entity.Word) (int64, error) {
	query := `DELETE FROM dictionary WHERE id=$1`
	result, err := r.db.ExecContext(ctx, query, w.ID)
	if err != nil {
		return 0, fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWord - exec: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dictionary.repository.DictionaryRepo.DeleteWord - rows affected: %w", err)
	}
	return rows, nil
}

func (r *DictionaryRepo) GetRandomWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	table := getTable(w.LangCode)
	query := fmt.Sprintf(`SELECT text, pronunciation, lang_code FROM "%s" WHERE moderator IS NOT NULL ORDER BY RANDOM() LIMIT 1;`, table)
	word := &entity.Word{}
	err := r.db.QueryRowContext(ctx, query).Scan(&word.Text, &word.Pronunciation, &word.LangCode)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictionaryRepo.GetRandomWord - scan: %w", err)
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
