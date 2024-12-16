package service

import (
	"context"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/games"

	"github.com/google/uuid"
)

type (
	repo interface{}
)

type Service struct {
	repo repo
}

func New(repo repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GameRevise(ctx context.Context, uid uuid.UUID, data entity.ReviseGame) error {
	return nil
}
