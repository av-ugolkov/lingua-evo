package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (d *Database) AddWordInDictionary(ctx context.Context, userId string, originalWord uuid.UUID, translateWord uuid.UUID) error {
	query := `INSERT INTO dictionary (user_id, original_word, original_lang, translate_word, translate_lang, example) VALUES($1, $2, $3, $4, $5, $6)`
	err := d.db.QueryRowContext(ctx, query, userId, "original_word", "lang", "translate_word", "translate_lang", "example") //TODO fix parameters
	if err != nil {
		return fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return nil
}

func (d *Database) GetWordsByUser(ctx context.Context, userId string) ([]Word, error) {
	return []Word{}, nil
}

func (d *Database) GetRandomWordByUser(ctx context.Context, userId string) (Word, error) {
	return Word{}, nil
}
