package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type UsersDB interface {
	AddUser(ctx context.Context, u *User) (uuid.UUID, error)
	EditUser(ctx context.Context, u *User) error
	FindUser(ctx context.Context, username string) (*User, error)
	RemoveUser(ctx context.Context, u *User) error
}

func (d *Database) AddUser(ctx context.Context, u *User) (uuid.UUID, error) {
	query := `INSERT INTO users (name, password_hash, last_visit) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id`

	var uid uuid.UUID

	err := d.db.QueryRowContext(ctx, query, u.Username, u.PasswordHash, time.Now()).Scan(&uid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("database.AddUser.QueryRow: %w", err)
	}

	return uid, nil
}

func (d *Database) EditUser(ctx context.Context, u *User) error {
	return nil
}

func (d *Database) FindUser(ctx context.Context, username string) (uuid.UUID, error) {
	query := `SELECT id FROM users where name=$1`

	var uid uuid.UUID

	err := d.db.QueryRowContext(ctx, query, username).Scan(&uid)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (d *Database) RemoveUser(ctx context.Context, u *User) error {
	return nil
}
