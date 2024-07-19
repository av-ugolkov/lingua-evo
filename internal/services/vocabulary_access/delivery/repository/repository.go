package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VocabAccessRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *VocabAccessRepo {
	return &VocabAccessRepo{
		pgxPool: pgxPool,
	}
}

func (r *VocabAccessRepo) AddAccessForUser(ctx context.Context, vocabID, userID uuid.UUID, isEditor bool) error {
	const query = `INSERT INTO vocabulary_users_access (vocab_id, subscriber_id, editor) VALUES ($1, $2, $3);`
	result, err := r.pgxPool.Exec(ctx, query, vocabID, userID, isEditor)
	if err != nil {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.AddAccess: %w", err)
	}

	if rows := result.RowsAffected(); rows != 1 {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.AddAccess: change 0 or more than 1 rows")
	}

	return nil
}

func (r *VocabAccessRepo) RemoveAccessForUser(ctx context.Context, vocabID, userID uuid.UUID) error {
	const query = `DELETE FROM vocabulary_users_access where vocab_id=$1 AND subscriber_id=$2;`
	result, err := r.pgxPool.Exec(ctx, query, vocabID, userID)
	if err != nil {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.RemoveAccess: %w", err)
	}

	if rows := result.RowsAffected(); rows != 1 {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.RemoveAccess: change 0 or more than 1 rows")
	}

	return nil
}

func (r *VocabAccessRepo) GetAccess(ctx context.Context, vocabID, userID uuid.UUID) (bool, error) {
	const query = `SELECT editor FROM vocabulary_users_access WHERE vocab_id=$1 AND subscriber_id=$2;`

	var isEditor bool
	err := r.pgxPool.QueryRow(ctx, query, vocabID, userID).Scan(&isEditor)
	if err != nil {
		return false, fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.GetAccess: %w", err)
	}

	return isEditor, nil
}
