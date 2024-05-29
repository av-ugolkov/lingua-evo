package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type VocabAccessRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *VocabAccessRepo {
	return &VocabAccessRepo{
		db: db,
	}
}

func (r *VocabAccessRepo) ChangeAccess(ctx context.Context, vocabID uuid.UUID, access int, accessEdit bool) error {
	const query = `UPDATE vocabulary SET access=$2, access_edit=$3 WHERE id=$1;`
	result, err := r.db.ExecContext(ctx, query, vocabID, access, accessEdit)
	if err != nil {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.ChangeAccess: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows != 1 {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.ChangeAccess: change 0 or more than 1 rows")
	}

	return nil
}

func (r *VocabAccessRepo) AddAccessForUser(ctx context.Context, vocabID, userID uuid.UUID, isEditor bool) error {
	const query = `INSERT INTO vocabulary_users_access (vocab_id, subscriber_id, editor) VALUES ($1, $2, $3);`
	result, err := r.db.ExecContext(ctx, query, vocabID, userID, isEditor)
	if err != nil {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.AddAccess: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows != 1 {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.AddAccess: change 0 or more than 1 rows")
	}

	return nil
}

func (r *VocabAccessRepo) RemoveAccessForUser(ctx context.Context, vocabID, userID uuid.UUID) error {
	const query = `DELETE FROM vocabulary_users_access where vocab_id=$1 AND subscriber_id=$2;`
	result, err := r.db.ExecContext(ctx, query, vocabID, userID)
	if err != nil {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.RemoveAccess: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows != 1 {
		return fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.RemoveAccess: change 0 or more than 1 rows")
	}

	return nil
}

func (r *VocabAccessRepo) GetAccess(ctx context.Context, vocabID, userID uuid.UUID) (bool, error) {
	const query = `SELECT editor FROM vocabulary_users_access WHERE vocab_id=$1 AND subscriber_id=$2;`

	var isEditor bool
	if err := r.db.QueryRowContext(ctx, query, vocabID, userID).Scan(&isEditor); err != nil {
		return false, fmt.Errorf("vocabulary_access.repository.VocabAccessRepo.GetAccess: %w", err)
	}

	return isEditor, nil
}
