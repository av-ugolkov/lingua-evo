package repository

import (
	"context"
	"database/sql"
	"fmt"

	entityWord "lingua-evo/internal/services/word/entity"

	"github.com/google/uuid"
)

type DictRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *DictRepo {
	return &DictRepo{
		db: db,
	}
}

func (r *DictRepo) AddWordInDictionary(ctx context.Context, userId, originalWord uuid.UUID, translateWord []uuid.UUID, pronunciation string, examples []uuid.UUID) error {
	query := `INSERT INTO dictionary (user_id, original_word, translate_word, pronunciation, examples) VALUES($1, $2, $3, $4, $5)`
	_, err := r.db.QueryContext(ctx, query, userId, originalWord, translateWord, pronunciation, examples)
	if err != nil {
		return fmt.Errorf("database.AddWord.QueryRow: %v", err)
	}

	return nil
}

func (r *DictRepo) GetWordsByUser(ctx context.Context, userId string) ([]entityWord.Word, error) {
	return []entityWord.Word{}, nil
}

func (r *DictRepo) GetRandomWordByUser(ctx context.Context, userId string) (entityWord.Word, error) {
	return entityWord.Word{}, nil
}
