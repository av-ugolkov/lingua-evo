package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	entity "lingua-evo/internal/services/user"
	"lingua-evo/runtime"
)

type (
	userRepo interface {
		AddUser(ctx context.Context, u *entity.User) (uuid.UUID, error)
		EditUser(ctx context.Context, u *entity.User) error
		GetUserByID(ctx context.Context, uid uuid.UUID) (*entity.User, error)
		GetUserByName(ctx context.Context, name string) (*entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
		GetUserByToken(ctx context.Context, token uuid.UUID) (*entity.User, error)
		RemoveUser(ctx context.Context, u *entity.User) error
	}

	redis interface {
		Get(ctx context.Context, key string) (string, error)
	}

	UserSvc struct {
		repo  userRepo
		redis redis
	}
)

func NewService(repo userRepo, redis redis) *UserSvc {
	return &UserSvc{
		repo:  repo,
		redis: redis,
	}
}

func (s *UserSvc) CreateUser(ctx context.Context, username, password, email string, role runtime.Role) (uuid.UUID, error) {
	user := &entity.User{
		ID:           uuid.New(),
		Name:         username,
		PasswordHash: password,
		Email:        email,
		Role:         role,
		CreatedAt:    time.Now().UTC(),
		LastVisitAt:  time.Now().UTC(),
	}

	uid, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (s *UserSvc) EditUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (s *UserSvc) GetUser(ctx context.Context, login string) (*entity.User, error) {
	user, err := s.repo.GetUserByName(ctx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user.service.UserSvc.GetUser - by name: %w", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		user, err = s.repo.GetUserByEmail(ctx, login)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user.service.UserSvc.GetUser - by email: %w", err)
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user.service.UserSvc.GetUser - by [%s]: %w", login, entity.ErrNotFoundUser)
	}

	return user, nil
}

func (s *UserSvc) GetUserByID(ctx context.Context, uid uuid.UUID) (*entity.User, error) {
	return s.repo.GetUserByID(ctx, uid)
}

func (s *UserSvc) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	return s.repo.GetUserByName(ctx, name)
}

func (s *UserSvc) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserSvc) GetUserByRefreshToken(ctx context.Context, token uuid.UUID) (*entity.User, error) {
	sessionJson, err := s.redis.Get(ctx, token.String())
	if err != nil {
		return nil, fmt.Errorf("user.service.UserSvc.GetUserByRefreshToken: %w", err)
	}

	var session entity.Session

	err = json.Unmarshal([]byte(sessionJson), &session)
	if err != nil {
		return nil, fmt.Errorf("user.service.UserSvc.GetUserByRefreshToken: %w", err)
	}

	return s.repo.GetUserByToken(ctx, session.UserID)
}

func (s *UserSvc) RemoveUser(ctx context.Context, user *entity.User) error {
	return nil
}
