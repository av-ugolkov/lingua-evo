package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"lingua-evo/internal/services/user/dto"
	"lingua-evo/internal/services/user/entity"
)

type userRepo interface {
	AddUser(ctx context.Context, u *entity.User) (uuid.UUID, error)
	EditUser(ctx context.Context, u *entity.User) error
	GetUserByName(ctx context.Context, name string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	RemoveUser(ctx context.Context, u *entity.User) error
}

type UserSvc struct {
	repo userRepo
}

func NewService(repo userRepo) *UserSvc {
	return &UserSvc{
		repo: repo,
	}
}

func (s *UserSvc) CreateUser(ctx context.Context, u *dto.CreateUserRq) (uuid.UUID, error) {
	user := &entity.User{
		ID:           uuid.New(),
		Username:     u.Username,
		PasswordHash: u.Password,
		Email:        u.Email,
		Role:         "user",
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

func (s *UserSvc) GetUserByName(ctx context.Context, name string) (*entity.User, error) {
	return s.repo.GetUserByName(ctx, name)
}

func (s *UserSvc) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *UserSvc) RemoveUser(ctx context.Context, user *entity.User) error {
	return nil
}
