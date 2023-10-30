package repository

import (
	"context"
	"database/sql"
	"fmt"
	"lingua-evo/internal/services/auth/entity"

	"github.com/google/uuid"
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
	query := `INSERT INTO session (refresh_token, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, s.RefreshToken, s.UserID, s.ExpiresAt, s.CreatedAt)
	if err != nil {
		return fmt.Errorf("auth.repository.SessionRepo.CreateSession: %w", err)
	}
	return nil
}

func (r *SessionRepo) GetSession(ctx context.Context, refreshToken uuid.UUID) (*entity.Session, error) {
	var s entity.Session
	query := `SELECT refresh_token, user_id, expires_at, created_at FROM session WHERE refresh_token=$1`
	err := r.db.QueryRowContext(ctx, query, refreshToken).Scan(&s.RefreshToken, &s.UserID, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("auth.repository.SessionRepo.GetSession: %w", err)
	}
	return &s, nil
}

func (r *SessionRepo) GetCountSession(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	query := `SELECT count(*) FROM session WHERE user_id=$1`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("auth.repository.SessionRepo.GetCountSession: %w", err)
	}
	return count, nil
}

func (r *SessionRepo) DeleteSession(ctx context.Context, session uuid.UUID) error {
	query := `DELETE FROM session WHERE refresh_token=$1`
	_, err := r.db.ExecContext(ctx, query, session)
	if err != nil {
		return fmt.Errorf("auth.repository.SessionRepo.DeleteSession: %w", err)
	}
	return nil
}

func (r *SessionRepo) DeleteAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM session WHERE user_id=$1`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("auth.repository.SessionRepo.DeleteAllUserSessions: %w", err)
	}
	return nil
}
