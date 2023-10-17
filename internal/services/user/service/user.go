package service

import (
	"context"

	"github.com/google/uuid"

	"lingua-evo/internal/services/user/entity"
)

type userRepo interface {
	AddUser(ctx context.Context, u *entity.User) (uuid.UUID, error)
	EditUser(ctx context.Context, u *entity.User) error
	GetIDByName(ctx context.Context, name string) (uuid.UUID, error)
	GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error)
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

func (s *UserSvc) CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	uid, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return uid, nil
}

func (s *UserSvc) EditUser(ctx context.Context, user *entity.User) error {
	return nil
}

func (s *UserSvc) GetIDByName(ctx context.Context, name string) (uuid.UUID, error) {
	return s.repo.GetIDByName(ctx, name)
}

func (s *UserSvc) GetIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	return s.repo.GetIDByEmail(ctx, email)
}

func (s *UserSvc) RemoveUser(ctx context.Context, user *entity.User) error {
	return nil
}
