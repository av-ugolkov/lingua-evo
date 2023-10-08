package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lingua-evo/internal/services/dictionary/entity"

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

func (r *DictRepo) AddDictionary(ctx context.Context, userID uuid.UUID, name string) (uuid.UUID, error) {
	query := `INSERT INTO dictionary (id, user_id, name) VALUES($1, $2, $3)`
	dictID := uuid.New()
	_, err := r.db.QueryContext(ctx, query, dictID, userID, name)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.repository.DictRepo.AddDictionary: %w", err)
	}
	return dictID, nil
}

func (r *DictRepo) DeleteDictionary(ctx context.Context, userID uuid.UUID, name string) error {
	query := `DELETE FROM dictionary WHERE user_id=$1 AND name=$2;`
	_, err := r.db.QueryContext(ctx, query, userID, name)
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictRepo.AddDictionary: %w", err)
	}
	return nil
}

func (r *DictRepo) GetDictionary(ctx context.Context, userID uuid.UUID, name string) (uuid.UUID, error) {
	query := `SELECT id FROM dictionary WHERE user_id=$1 AND name=$2;`
	_, err := r.db.QueryContext(ctx, query, userID, name)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.repository.DictRepo.AddDictionary: %w", err)
	}
	return uuid.Nil, nil
}

func (r *DictRepo) GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*entity.Dictionary, error) {
	query := `SELECT id, user_id, name FROM dictionary WHERE user_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("dictionary.repository.DictRepo.AddDictionary: %w", err)
	}
	var dictionaries []*entity.Dictionary
	for rows.Next() {
		var dict entity.Dictionary
		err := rows.Scan(
			&dict.ID,
			&dict.UserID,
			&dict.Name,
		)
		if err != nil {
			return nil, err
		}

		dictionaries = append(dictionaries, &dict)
	}

	return dictionaries, nil
}
