package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (d *Database) AddWordInDictionary(ctx context.Context, userId, originalWord uuid.UUID, translateWord []uuid.UUID, pronunciation string, examples []uuid.UUID) error {
	query := `INSERT INTO dictionary (user_id, original_word, translate_word, pronunciation, examples) VALUES($1, $2, $3, $4, $5)`
	_, err := d.db.QueryContext(ctx, query, userId, originalWord, translateWord, pronunciation, examples)
	if err != nil {
		return fmt.Errorf("database.AddWord.QueryRow: %v", err)
	}

	return nil
}

func (d *Database) GetWordsByUser(ctx context.Context, userId string) ([]Word, error) {
	return []Word{}, nil
}

func (d *Database) GetRandomWordByUser(ctx context.Context, userId string) (Word, error) {
	return Word{}, nil
}
