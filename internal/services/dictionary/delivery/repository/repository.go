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

func (r *DictRepo) Add(ctx context.Context, dict entity.Dictionary) error {
	query := `INSERT INTO dictionary (id, user_id, name, native_lang_code, second_lang_code, tags) VALUES($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query, dict.ID, dict.UserID, dict.Name, dict.NativeLang, dict.SecondLang, dict.Tags)
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictRepo.AddDictionary: %w", err)
	}

	return nil
}

func (r *DictRepo) Delete(ctx context.Context, dict entity.Dictionary) error {
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

func (r *DictRepo) GetByName(ctx context.Context, userID uuid.UUID, name string) (entity.Dictionary, error) {
	query := `SELECT id, native_lang_code, second_lang_code, tags FROM dictionary WHERE user_id=$1 AND name=$2;`
	var dict entity.Dictionary
	err := r.db.QueryRowContext(ctx, query, userID, name).Scan(&dict.ID, &dict.NativeLang, &dict.SecondLang, pq.Array(&dict.Tags))
	if err != nil {
		return dict, err
	}
	return dict, nil
}

func (r *DictRepo) GetByID(ctx context.Context, dictID uuid.UUID) (entity.Dictionary, error) {
	query := `SELECT user_id, name, native_lang_code, second_lang_code, tags FROM dictionary WHERE id=$1;`
	var dict entity.Dictionary
	err := r.db.QueryRowContext(ctx, query, dictID).Scan(&dict.UserID, &dict.Name, &dict.NativeLang, &dict.SecondLang, pq.Array(&dict.Tags))
	if err != nil {
		return dict, err
	}
	return dict, nil
}

func (r *DictRepo) GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*entity.Dictionary, error) {
	query := `SELECT d.id, d.user_id, name, n.lang native_lang_code, s.lang second_lang_code FROM dictionary d
left join "language" n on n.code = d.native_lang_code
left join "language" s on s.code = d.second_lang_code 
WHERE user_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dictionaries []*entity.Dictionary
	for rows.Next() {
		var dict entity.Dictionary
		err := rows.Scan(
			&dict.ID,
			&dict.UserID,
			&dict.Name,
			&dict.NativeLang,
			&dict.SecondLang,
		)
		if err != nil {
			return nil, fmt.Errorf("dictionary.repository.DictRepo.GetDictionaries - scan: %w", err)
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

func (r *DictRepo) Rename(ctx context.Context, id uuid.UUID, newName string) error {
	query := `UPDATE dictionary SET name=$1 WHERE id=$2;`
	result, err := r.db.ExecContext(ctx, query, newName, id)
	if err != nil {
		return fmt.Errorf("dictionary.repository.DictRepo.RenameDictionary: %w", err)
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("dictionary.repository.DictRepo.RenameDictionary: %w", entity.ErrDictionaryNotFound)
	}
	return nil
}
