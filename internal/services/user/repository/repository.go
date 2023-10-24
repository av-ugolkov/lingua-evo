package repository

import (
	"context"
	"database/sql"
	"fmt"

	"lingua-evo/internal/services/user/entity"

	"github.com/google/uuid"
)

type UserRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) AddUser(ctx context.Context, u *entity.User) (uuid.UUID, error) {
	query := `INSERT INTO users (id, name, email, password_hash, role, last_visit_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var uid uuid.UUID

	err := r.db.QueryRowContext(ctx, query, u.ID, u.Username, u.Email, u.PasswordHash, u.Role, u.LastVisitAt, u.CreatedAt).Scan(&uid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user.repository.AddUser: %w", err)
	}

	return uid, nil
}

func (r *UserRepo) EditUser(ctx context.Context, u *entity.User) error {
	return nil
}

func (r *UserRepo) GetIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	query := `SELECT id FROM users where name=$1`

	var uid uuid.UUID

	err := r.db.QueryRowContext(ctx, query, name).Scan(&uid)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (r *UserRepo) GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	query := `SELECT id FROM users where email=$1`

	var uid uuid.UUID

	err := r.db.QueryRowContext(ctx, query, email).Scan(&uid)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (r *UserRepo) RemoveUser(ctx context.Context, u *entity.User) error {
	return nil
}
