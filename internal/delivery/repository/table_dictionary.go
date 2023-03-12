package repository

import (
	"context"
	"fmt"
)

func (d *Database) AddWordInDictionary(ctx context.Context, userId string, w *Word) error {
	query := `insert into dictionary (user_id, original_word, original_lang, translate_word, translate_lang, example) values($1, $2, $3, $4, $5, $6)`
	err := d.db.QueryRowContext(ctx, query, userId, "original_word", "lang", "translate_word", "translate_lang", "example") //TODO fix parameters
	if err != nil {
		return fmt.Errorf("database.AddWord.QueryRow: %w", err)
	}

	return nil
}
