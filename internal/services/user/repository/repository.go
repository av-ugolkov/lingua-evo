package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	sorted "github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

type UserRepo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *UserRepo {
	return &UserRepo{
		tr: tr,
	}
}

func (r *UserRepo) AddUser(ctx context.Context, u *entity.User, pswHash string) (uuid.UUID, error) {
	query := `
		INSERT INTO users (
			id, 
			nickname, 
			email, 
			password_hash, 
			role, 
			visited_at, 
			created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var uid uuid.UUID
	now := time.Now().UTC()
	err := r.tr.QueryRow(ctx, query, uuid.New(), u.Nickname, u.Email, pswHash, u.Role, now, now).Scan(&uid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user.repository.UserRepo.AddUser: %w", err)
	}

	return uid, nil
}

func (r *UserRepo) AddGoogleUser(ctx context.Context, u entity.GoogleUser) (uuid.UUID, error) {
	query := `
		INSERT INTO users (
			id, 
			nickname, 
			email, 
			role, 
			google_id,
			visited_at, 
			created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var uid uuid.UUID
	now := time.Now().UTC()
	err := r.tr.QueryRow(ctx, query, uuid.New(), u.Nickname, u.Email, u.Role, u.GoogleID, now, now).Scan(&uid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user.repository.UserRepo.AddGoogleUser: %w", err)
	}

	return uid, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, uid uuid.UUID) (*entity.User, error) {
	query := `
		SELECT 
			id, 
			nickname, 
			email, 
			role, 
			visited_at, 
			created_at 
		FROM users WHERE id=$1`
	var u entity.User
	err := r.tr.QueryRow(ctx, query, uid).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
		&u.Role,
		&u.VisitedAt,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByID: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetUserByNickname(ctx context.Context, nickname string) (*entity.User, error) {
	query := `
		SELECT 
			id, 
			nickname, 
			email, 
			role,
			visited_at,
			created_at 
		FROM users WHERE nickname=$1`

	var u entity.User
	err := r.tr.QueryRow(ctx, query, nickname).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
		&u.Role,
		&u.VisitedAt,
		&u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByNickname: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT 
			id,
			nickname,
			email,
			role,
			visited_at,
			created_at 
		FROM users WHERE email=$1`

	var u entity.User

	err := r.tr.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
		&u.Role,
		&u.VisitedAt,
		&u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByEmail: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) GetUserByGoogleID(ctx context.Context, googleID string) (*entity.User, error) {
	query := `
		SELECT 
			id,
			nickname,
			email,
			role,
			visited_at,
			created_at 
		FROM users WHERE google_id=$1`

	var u entity.User

	err := r.tr.QueryRow(ctx, query, googleID).Scan(
		&u.ID,
		&u.Nickname,
		&u.Email,
		&u.Role,
		&u.VisitedAt,
		&u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByGoogleID: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) RemoveUser(ctx context.Context, uid uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.tr.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.RemoveUser: %w", err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("user.repository.UserRepo.RemoveUser: change 0 or more than 1 rows")
	}

	return nil
}

func (r *UserRepo) GetUserData(ctx context.Context, uid uuid.UUID) (*entity.UserData, error) {
	const query = `SELECT id, name, surname FROM user_data WHERE user_id = $1`

	var data entity.UserData
	err := r.tr.QueryRow(ctx, query, uid).Scan(&data.UID, &data.Name, &data.Surname)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserData: %w", err)
	}

	return &data, nil
}

func (r *UserRepo) AddUserData(ctx context.Context, data entity.UserData) error {
	const query = `INSERT INTO user_data(user_id, name, surname) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	_, err := r.tr.Exec(ctx, query, data.UID, data.Name, data.Surname)
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.GetUserData: %w", err)
	}

	return nil
}

func (r *UserRepo) AddUserNewsletters(ctx context.Context, data entity.UserNewsletters) error {
	const query = `INSERT INTO user_newsletters(user_id, news) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.tr.Exec(ctx, query, data.UID, data.News)
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.AddUserNewsletters: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUserSubscriptions(ctx context.Context, uid uuid.UUID) ([]entity.Subscriptions, error) {
	const query = `
	SELECT us.id, user_id, subscription_id, s.add_words count_word, us.started_at, us.ended_at 
	FROM user_subscription us
	LEFT JOIN subscriptions s ON s.id = us.subscription_id
	WHERE user_id = $1;`

	rows, err := r.tr.Query(ctx, query, uid)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserSubscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []entity.Subscriptions
	for rows.Next() {
		var sub entity.Subscriptions
		if err := rows.Scan(&sub.ID, &sub.UserID, &sub.SubscriptionID, &sub.CountWord, &sub.StartedAt, &sub.EndedAt); err != nil {
			return nil, fmt.Errorf("user.repository.UserRepo.GetUserSubscriptions: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *UserRepo) GetUsers(ctx context.Context, page, perPage, sort, order int, search string) ([]entity.User, int, error) {
	const queryCountUsers = `SELECT COUNT(id) FROM users`
	var countUser int
	err := r.tr.QueryRow(ctx, queryCountUsers).Scan(&countUser)
	if err != nil {
		return nil, 0, fmt.Errorf("user.repository.UserRepo.GetUsers: %w", err)
	}

	query := fmt.Sprintf(`
	SELECT u.id, u.nickname, role, u.visited_at
	FROM users u
	WHERE u.nickname LIKE '%[1]s' || $1 || '%[1]s'
	%[2]s
	LIMIT $2
	OFFSET $3;`, "%", getSorted(sort, sorted.TypeOrder(order)))

	rows, err := r.tr.Query(ctx, query, search, perPage, (page-1)*perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("user.repository.UserRepo.GetUsers: %w", err)
	}
	defer rows.Close()

	users := make([]entity.User, 0)
	var user entity.User
	for rows.Next() {
		if err := rows.Scan(
			&user.ID,
			&user.Nickname,
			&user.Role,
			&user.VisitedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("user.repository.UserRepo.GetUsers: %w", err)
		}
		users = append(users, user)
	}

	return users, countUser, nil
}

func (r *UserRepo) UpdateVisitedAt(ctx context.Context, uid uuid.UUID) error {
	const query = `UPDATE users SET visited_at = $2 WHERE id = $1`

	_, err := r.tr.Exec(ctx, query, uid, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.UpdateLastVisited: %w", err)
	}

	return nil
}

func getSorted(typeSorted int, order sorted.TypeOrder) string {
	switch sorted.TypeSorted(typeSorted) {
	case sorted.Visit:
		return fmt.Sprintf("ORDER BY u.visited_at %s", order)
	case sorted.ABC:
		return fmt.Sprintf("ORDER BY u.nickname %s", order)
	default:
		return runtime.EmptyString
	}
}
