package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	sorted "github.com/av-ugolkov/lingua-evo/internal/pkg/utils"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/user"

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

func (r *UserRepo) AddUser(ctx context.Context, u *entity.User) (uuid.UUID, error) {
	query := `INSERT INTO users (id, name, email, password_hash, role, last_visit_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING id`

	var uid uuid.UUID
	err := r.tr.QueryRow(ctx, query, uuid.New(), u.Name, u.Email, u.PasswordHash, u.Role, u.LastVisitAt, u.CreatedAt).Scan(&uid)
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
	err := r.tr.QueryRow(ctx, query, uid).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByID: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	query := `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users where name=$1`

	var u entity.User

	err := r.tr.QueryRow(ctx, query, name).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByName: %w", err)
	}

	return &u, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, user, email, password_hash, role, last_visit_at, created_at FROM users where email=$1`

	var u entity.User

	err := r.tr.QueryRow(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByEmail: %w", err)
	}

	u.Email = email

	return &u, nil
}

func (r *UserRepo) GetUserByToken(ctx context.Context, token uuid.UUID) (*entity.User, error) {
	query := `SELECT id, name, email, password_hash, role, last_visit_at, created_at FROM users where id = $1`

	var u entity.User
	err := r.tr.QueryRow(ctx, query, token).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.LastVisitAt, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserByToken: %w", err)
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

func (r *UserRepo) GetUserData(ctx context.Context, userID uuid.UUID) (*entity.Data, error) {
	const query = `SELECT max_count_words, newsletter FROM user_data WHERE user_id = $1`

	var data entity.Data
	err := r.tr.QueryRow(ctx, query, userID).Scan(&data.MaxCountWords, &data.Newsletters)
	if err != nil {
		return nil, fmt.Errorf("user.repository.UserRepo.GetUserData: %w", err)
	}

	return &data, nil
}

func (r *UserRepo) AddUserData(ctx context.Context, userID uuid.UUID, maxCountWords int, newsletter bool) error {
	const query = `INSERT INTO user_data(user_id, max_count_words, newsletter) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	_, err := r.tr.Exec(ctx, query, userID, maxCountWords, newsletter)
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.GetUserData: %w", err)
	}

	return nil
}

func (r *UserRepo) GetUserSubscriptions(ctx context.Context, userID uuid.UUID) ([]entity.Subscriptions, error) {
	const query = `
	SELECT us.id, user_id, subscription_id, s.add_words count_word, us.started_at, us.ended_at 
	FROM user_subscription us
	LEFT JOIN subscriptions s ON s.id = us.subscription_id
	WHERE user_id = $1;`

	rows, err := r.tr.Query(ctx, query, userID)
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

func (r *UserRepo) GetUsers(ctx context.Context, page, perPage, sort, order int, search string) ([]entity.UserData, int, error) {
	const queryCountUsers = `SELECT COUNT(id) FROM users`
	var countUser int
	err := r.tr.QueryRow(ctx, queryCountUsers).Scan(&countUser)
	if err != nil {
		return nil, 0, fmt.Errorf("user.repository.UserRepo.GetUsers: %w", err)
	}

	query := fmt.Sprintf(`
	SELECT u.id, u.name, role, u.last_visit_at
	FROM users u
	WHERE POSITION($1 in u."name")>0
	%s
	LIMIT $2
	OFFSET $3;`, getSorted(sort, sorted.TypeOrder(order)))

	rows, err := r.tr.Query(ctx, query, search, perPage, (page-1)*perPage)
	if err != nil {
		return nil, 0, fmt.Errorf("user.repository.UserRepo.GetUsers: %w", err)
	}
	defer rows.Close()

	users := make([]entity.UserData, 0)
	var user entity.UserData
	for rows.Next() {
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Role,
			&user.LastVisited,
		); err != nil {
			return nil, 0, fmt.Errorf("user.repository.UserRepo.GetUsers - scan: %w", err)
		}
		users = append(users, user)
	}

	return users, countUser, nil
}

func (r *UserRepo) UpdateLastVisited(ctx context.Context, uid uuid.UUID) error {
	const query = `UPDATE users SET last_visit_at = $2 WHERE id = $1`

	_, err := r.tr.Exec(ctx, query, uid, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("user.repository.UserRepo.UpdateLastVisited: %w", err)
	}

	return nil
}

func getSorted(typeSorted int, order sorted.TypeOrder) string {
	switch sorted.TypeSorted(typeSorted) {
	case sorted.Visit:
		return fmt.Sprintf("ORDER BY u.last_visit_at %s", order)
	case sorted.ABC:
		return fmt.Sprintf("ORDER BY u.name %s", order)
	default:
		return ""
	}
}
