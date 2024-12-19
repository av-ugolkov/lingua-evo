package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"

	"github.com/google/uuid"
)

type Repo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *Repo {
	return &Repo{
		tr: tr,
	}
}

func (r *Repo) Get(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT subscribers_id FROM subscribers WHERE user_id=$1;`
	rows, err := r.tr.Query(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("subscribers.repository.SubscribersRepo.Get: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("subscribers.repository.SubscribersRepo.Get: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repo) GetRespondents(ctx context.Context, uid uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM subscribers WHERE subscribers_id=$1;`
	rows, err := r.tr.Query(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("subscribers.repository.SubscribersRepo.GetRespondents: %w", err)
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("subscribers.repository.SubscribersRepo.GetRespondents: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *Repo) Subscribe(ctx context.Context, uid, subID uuid.UUID) error {
	const query = `INSERT INTO subscribers (user_id, subscribers_id, created_at) VALUES ($1, $2, $3);`

	_, err := r.tr.Exec(ctx, query, uid, subID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("subscribers.repository.SubscribersRepo.Subscribe: %w", err)
	}

	return nil
}

func (r *Repo) Unsubscribe(ctx context.Context, uid, subID uuid.UUID) error {
	const query = `DELETE FROM subscribers WHERE user_id=$1 AND subscribers_id=$2;`

	result, err := r.tr.Exec(ctx, query, uid, subID)
	if err != nil {
		return fmt.Errorf("subscribers.repository.SubscribersRepo.Unsubscribe: %w", err)
	}

	if rows := result.RowsAffected(); rows != 1 {
		return fmt.Errorf("subscribers.repository.SubscribersRepo.Unsubscribe: change 0 or more than 1 rows")
	}

	return nil
}

func (r *Repo) Check(ctx context.Context, uid, subID uuid.UUID) (bool, error) {
	query := `SELECT count(user_id) FROM subscribers WHERE user_id=$1 AND subscribers_id=$2;`

	var count int
	err := r.tr.QueryRow(ctx, query, uid, subID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("subscribers.repository.SubscribersRepo.Get: %w", err)
	}

	return count != 0, nil
}
