package repository

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/access"
)

type AccessRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *AccessRepo {
	return &AccessRepo{
		db: db,
	}
}

func (r *AccessRepo) GetAccesses(ctx context.Context) ([]entity.Access, error) {
	query := `SELECT id, type, name FROM access`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("access.delivery.repository.AccessRepo.GetAccesses: %w", err)
	}
	defer rows.Close()

	accesses := make([]entity.Access, 0)
	for rows.Next() {
		var access entity.Access
		if err := rows.Scan(&access.ID, &access.Type, &access.Name); err != nil {
			return nil, fmt.Errorf("access.delivery.repository.AccessRepo.GetAccesses: %w", err)
		}
		accesses = append(accesses, access)
	}
	return accesses, nil

}
