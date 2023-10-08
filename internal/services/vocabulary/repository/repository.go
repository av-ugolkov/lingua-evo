package repository

import (
	"context"
	"database/sql"
	"fmt"

	entityWord "lingua-evo/internal/services/word/entity"

	"github.com/google/uuid"
)

type VocabularyRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *VocabularyRepo {
	return &VocabularyRepo{
		db: db,
	}
}

func (r *VocabularyRepo) AddWordInVocabulary(ctx context.Context, userId, originalWord uuid.UUID, translateWord []uuid.UUID, pronunciation string, examples []uuid.UUID) error {
	//TODO переписать
	query := `INSERT INTO vocabulary (user_id, original_word, translate_word, pronunciation, examples) VALUES($1, $2, $3, $4, $5)`
	_, err := r.db.QueryContext(ctx, query, userId, originalWord, translateWord, pronunciation, examples)
	if err != nil {
		return fmt.Errorf("vocabulary.VocabularyRepo.AddWordInVocabulary: %w", err)
	}

	return nil
}

func (r *VocabularyRepo) GetWordsByUser(ctx context.Context, userId string) ([]entityWord.Word, error) {
	return []entityWord.Word{}, nil
}

func (r *VocabularyRepo) GetRandomWordByUser(ctx context.Context, userId string) (entityWord.Word, error) {
	return entityWord.Word{}, nil
}
