package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscribersRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *SubscribersRepo {
	return &SubscribersRepo{
		pgxPool: pgxPool,
	}
}

func (r *SubscribersRepo) Get(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT subscribers_id FROM subscribers WHERE user_id=$1;`
	rows, err := r.pgxPool.Query(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("subscribers.delivery.repository.SubscribersRepo.Get: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("subscribers.delivery.repository.SubscribersRepo.Get: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *SubscribersRepo) Check(ctx context.Context, uid, subID uuid.UUID) (bool, error) {
	query := `SELECT count(user_id) FROM subscribers WHERE user_id=$1 AND subscribers_id=$2;`

	var count int
	err := r.pgxPool.QueryRow(ctx, query, uid, subID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("subscribers.delivery.repository.SubscribersRepo.Get: %w", err)
	}

	return count != 0, nil
}
