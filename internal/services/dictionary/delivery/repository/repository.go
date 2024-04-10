package repository

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *DictionaryRepo) AddWord(ctx context.Context, w *entity.Word) (uuid.UUID, error) {
	var id uuid.UUID
	table := getTable(w.LangCode)
	query := fmt.Sprintf(
		`WITH d AS (
    		SELECT id FROM %[1]s WHERE text = $2),
		ins AS (
    		INSERT INTO %[1]s (id, text, pronunciation, lang_code, updated_at, created_at)
			VALUES($1, $2, $3, $4, $5, $5)
    		ON CONFLICT DO NOTHING RETURNING id)
		SELECT id
		FROM ins
		UNION ALL
		SELECT id
		FROM d;`, table)
	err := r.db.QueryRowContext(ctx, query, w.ID, w.Text, w.Pronunciation, w.LangCode, time.Now().UTC()).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.repository.DictionaryRepo.AddWord - query: %w", err)
	}

	return id, nil
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
