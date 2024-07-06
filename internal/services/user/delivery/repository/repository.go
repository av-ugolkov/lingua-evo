package repository

import (
	"context"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pgxPool: pgxPool,
	}
}

func (r *UserRepo) AddUser(ctx context.Context, u *entity.User) (uuid.UUID, error) {
	query := `INSERT INTO users (id, name, email, password_hash, role, last_visit_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var uid uuid.UUID
	err := r.pgxPool.QueryRow(ctx, query, u.ID, u.Name, u.Email, u.PasswordHash, u.Role, u.LastVisitAt, u.CreatedAt).Scan(&uid)
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
	err := r.pgxPool.QueryRow(ctx, query, uid).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByID: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	query := `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users where name=$1`

	var u entity.User

	err := r.pgxPool.QueryRow(ctx, query, name).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByName: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, user, email, password_hash, role, last_visit_at, created_at FROM users where email=$1`

	var u entity.User

	err := r.pgxPool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByEmail: %w", err)
	}

	u.Email = email

	return &u, nil
}

func (r *UserRepo) GetUserByToken(ctx context.Context, token uuid.UUID) (*entity.User, error) {
	query := `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users where id = $1`

	var u entity.User
	err := r.pgxPool.QueryRow(ctx, query, token).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByToken: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) RemoveUser(ctx context.Context, u *entity.User) error {
	return nil
}

func (r *UserRepo) GetUserData(ctx context.Context, userID uuid.UUID) (*entity.Data, error) {
	const query = `SELECT max_count_words, newsletter FROM user_data WHERE user_id = $1`

	var data entity.Data
	err := r.pgxPool.QueryRow(ctx, query, userID).Scan(&data.MaxCountWords, &data.Newsletters)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserData: %w", err)
	}

	return &data, nil
}

func (r *UserRepo) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]entity.Subscriptions, error) {
	const query = `
	SELECT us.id, user_id, subscription_id, s.add_words count_word, us.started_at, us.ended_at 
	FROM user_subscription us
	LEFT JOIN subscriptions s ON s.id = us.subscription_id
	WHERE user_id = $1;`

	rows, err := r.pgxPool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserSubscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []entity.Subscriptions
	for rows.Next() {
		var sub entity.Subscriptions
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.SubscriptionID, &sub.CountWord, &sub.StartedAt, &sub.EndedAt); err != nil {
			return nil, fmt.Errorf("user.repository.UserRepo.GetUserSubscriptions - scan: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}
