package repository

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DictRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *DictRepo {
	return &DictRepo{
		db: db,
	}
}

func (r *DictRepo) AddDictionary(ctx context.Context, dict entity.Dictionary) error {
	query := `INSERT INTO dictionary (id, user_id, name, tags) VALUES($1, $2, $3, $4)`

	_, err := r.db.QueryContext(ctx, query, dict.ID, dict.UserID, dict.Name, []uuid.UUID{uuid.New(), uuid.New()})
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictRepo.AddDictionary: %w", err)
	}
	return nil
}

func (r *DictRepo) DeleteDictionary(ctx context.Context, dict entity.Dictionary) error {
	query := `DELETE FROM dictionary WHERE user_id=$1 AND name=$2;`
	result, err := r.db.ExecContext(ctx, query, dict.UserID, dict.Name)
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictRepo.DeleteDictionary: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("dictionary.repository.DictRepo.DeleteDictionary: %w", entity.ErrDictionaryNotFound)
	}

	return nil
}

func (r *DictRepo) GetDictionaryByName(ctx context.Context, dict entity.Dictionary) (uuid.UUID, []uuid.UUID, error) {
	query := `SELECT id, tags FROM dictionary WHERE user_id=$1 AND name=$2;`
	var dictID uuid.UUID
	var tags []uuid.UUID
	err := r.db.QueryRowContext(ctx, query, dict.UserID, dict.Name).Scan(&dictID, pq.Array(&tags))
	if err != nil {
		return uuid.Nil, nil, err
	}
	return dictID, tags, nil
}

func (r *DictRepo) GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*entity.Dictionary, error) {
	query := `SELECT id, user_id, name FROM dictionary WHERE user_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
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

func (r *DictRepo) GetCountDictionaries(ctx context.Context, userID uuid.UUID) (int, error) {
	var countDictionaries int

	query := `SELECT COUNT(id) FROM dictionary WHERE user_id=$1;`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&countDictionaries)
	if err != nil {
		return -1, err
	}
	return countDictionaries, nil
}
