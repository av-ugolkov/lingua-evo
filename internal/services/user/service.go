package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/av-ugolkov/lingua-evo/runtime"
)

type (
	userRepo interface {
		AddUser(ctx context.Context, u *User) (uuid.UUID, error)
		EditUser(ctx context.Context, u *User) error
		GetUserByID(ctx context.Context, uid uuid.UUID) (*User, error)
		GetUserByName(ctx context.Context, name string) (*User, error)
		GetUserByEmail(ctx context.Context, email string) (*User, error)
		GetUserByToken(ctx context.Context, token uuid.UUID) (*User, error)
		RemoveUser(ctx context.Context, u *User) error
	}

	redis interface {
		Get(ctx context.Context, key string) (string, error)
	}

	Service struct {
		repo  userRepo
		redis redis
	}
)

func NewService(repo userRepo, redis redis) *Service {
	return &Service{
		repo:  repo,
		redis: redis,
	}
}

func (s *Service) CreateUser(ctx context.Context, username, password, email string, role runtime.Role) (uuid.UUID, error) {
	user := &User{
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

func (s *Service) EditUser(ctx context.Context, user *User) error {
	return nil
}

func (s *Service) GetUser(ctx context.Context, login string) (*User, error) {
	user, err := s.repo.GetUserByName(ctx, login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by name: %w", err)
	} else if errors.Is(err, sql.ErrNoRows) {
		user, err = s.repo.GetUserByEmail(ctx, login)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user.Service.GetUser - by email: %w", err)
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user.Service.GetUser - by [%s]: %w", login, ErrNotFoundUser)
	}

	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, uid uuid.UUID) (*User, error) {
	return s.repo.GetUserByID(ctx, uid)
}

func (s *Service) GetUserByName(ctx context.Context, name string) (*User, error) {
	return s.repo.GetUserByName(ctx, name)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *Service) GetUserByRefreshToken(ctx context.Context, token uuid.UUID) (*User, error) {
	sessionJson, err := s.redis.Get(ctx, token.String())
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	var session Session

	err = json.Unmarshal([]byte(sessionJson), &session)
	if err != nil {
		return nil, fmt.Errorf("user.Service.GetUserByRefreshToken: %w", err)
	}

	return s.repo.GetUserByToken(ctx, session.UserID)
}

func (s *Service) RemoveUser(ctx context.Context, user *User) error {
	return nil
}
