package repository

import (
	"context"
	"database/sql"
	"fmt"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"

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
	err := r.db.QueryRowContext(ctx, query, u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.LastVisitAt, u.CreatedAt).Scan(&uid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user.repository.UserRepo.AddUser: %w", err)
	}

	return uid, nil
}

func (r *UserRepo) EditUser(ctx context.Context, u *entity.User) error {
	return nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, uid uuid.UUID) (*entity.User, error) {
	query := `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users where id=$1`
	var u entity.User
	err := r.db.QueryRowContext(ctx, query, uid).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByID: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	query := `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users where name=$1`

	var u entity.User

	err := r.db.QueryRowContext(ctx, query, name).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByName: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, user, email, password_hash, role, last_visit_at, created_at FROM users where email=$1`

	var u entity.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByEmail: %w", err)
	}

	u.Email = email

	return &u, nil
}

func (r *UserRepo) GetUserByToken(ctx context.Context, token uuid.UUID) (*entity.User, error) {
	query := `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users where id = $1`

	var u entity.User
	err := r.db.QueryRowContext(ctx, query, token).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByToken: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) RemoveUser(ctx context.Context, u *entity.User) error {
	return nil
}
