package repository

import (
	"context"
	"fmt"
	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/google/uuid"
	"time"
)

type Repo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *Repo {
	return &Repo{tr: tr}
}

func (r *Repo) GetVocabNotification(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	const query = `SELECT user_id, vocab_id FROM vocabulary_notifications WHERE user_id = $1 AND vocab_id = $2`

	result, err := r.tr.Exec(ctx, query, uid, vid)
	if err != nil {
		return false, fmt.Errorf("notifications.delivery.repository.GetVocabNotifications: %w", err)
	}

	return result.RowsAffected() == 1, nil
}

func (r *Repo) SetVocabNotification(ctx context.Context, uid, vid uuid.UUID) error {
	const query = `INSERT INTO vocabulary_notifications(user_id, vocab_id, created_at)  VALUES ($1, $2, $3);`

	_, err := r.tr.Exec(ctx, query, uid, vid, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("notifications.delivery.repository.SetVocabNotifications: %w", err)
	}

	return nil
}

func (r *Repo) DeleteVocabNotification(ctx context.Context, uid, vid uuid.UUID) error {
	const query = `DELETE FROM vocabulary_notifications WHERE user_id=$1 AND vocab_id=$2;`

	_, err := r.tr.Exec(ctx, query, uid, vid)
	if err != nil {
		return fmt.Errorf("notifications.delivery.repository.DeleteVocabNotification: %w", err)
	}

	return nil
}
