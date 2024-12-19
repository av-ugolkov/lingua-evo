package repository

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/access"
)

type AccessRepo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *AccessRepo {
	return &AccessRepo{
		tr: tr,
	}
}

func (r *AccessRepo) GetAccesses(ctx context.Context) ([]entity.Access, error) {
	query := `SELECT id, type, name FROM access`
	rows, err := r.tr.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("access.repository.AccessRepo.GetAccesses: %w", err)
	}
	defer rows.Close()

	accesses := make([]entity.Access, 0)
	for rows.Next() {
		var access entity.Access
		if err := rows.Scan(&access.ID, &access.Type, &access.Name); err != nil {
			return nil, fmt.Errorf("access.repository.AccessRepo.GetAccesses: %w", err)
		}
		accesses = append(accesses, access)
	}
	return accesses, nil

}
