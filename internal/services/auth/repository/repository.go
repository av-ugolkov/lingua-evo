package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lingua-evo/internal/services/auth/entity"
)

type SessionRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{
		db: db,
	}
}

func (r *SessionRepo) SetSession(ctx context.Context, s *entity.Session) error {
	query := `INSERT INTO session (id, user_id, refresh_token, expires_at, created_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.UserID, s.RefreshToken, s.ExpiresAt, s.CreatedAt)
	if err != nil {
		return fmt.Errorf("auths.repository.SessionRepo.CreateSession: %w", err)
	}
	return nil
}
