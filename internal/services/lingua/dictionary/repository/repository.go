package repository

import (
	"context"
	"database/sql"

	"lingua-evo/internal/services/lingua/dictionary/entity"

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

func (r *DictRepo) AddDictionary(ctx context.Context, dict entity.Dictionary) error {
	query := `INSERT INTO dictionary (id, user_id, name) VALUES($1, $2, $3)`

	_, err := r.db.QueryContext(ctx, query, dict.ID, dict.UserID, dict.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *DictRepo) DeleteDictionary(ctx context.Context, dict entity.Dictionary) error {
	query := `DELETE FROM dictionary WHERE user_id=$1 AND name=$2;`
	_, err := r.db.QueryContext(ctx, query, dict.UserID, dict.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *DictRepo) GetDictionary(ctx context.Context, dict entity.Dictionary) (uuid.UUID, error) {
	query := `SELECT id FROM dictionary WHERE user_id=$1 AND name=$2;`
	var dictID uuid.UUID
	err := r.db.QueryRowContext(ctx, query, dict.UserID, dict.Name).Scan(&dictID)
	if err != nil {
		return uuid.Nil, err
	}
	return dictID, nil
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