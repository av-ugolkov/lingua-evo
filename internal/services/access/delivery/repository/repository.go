package repository

import (
	"context"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/access"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccessRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *AccessRepo {
	return &AccessRepo{
		pgxPool: pgxPool,
	}
}

func (r *AccessRepo) GetAccesses(ctx context.Context) ([]entity.Access, error) {
	query := `SELECT id, type, name FROM access`
	rows, err := r.pgxPool.Query(ctx, query)
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
