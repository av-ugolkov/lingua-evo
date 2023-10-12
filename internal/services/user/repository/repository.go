package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
	query := `INSERT INTO users (name, email, password_hash, last_visit) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING id`

	var uid uuid.UUID

	err := r.db.QueryRowContext(ctx, query, u.Username, u.Email, u.PasswordHash, time.Now()).Scan(&uid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("database.AddUser.QueryRow: %w", err)
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
