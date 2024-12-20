package service

import (
	"context"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/access"
)

type (
	repoAccess interface {
		GetAccesses(ctx context.Context) ([]entity.Access, error)
	}
)

type Service struct {
	repo repoAccess
}

func NewService(repo repoAccess) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetAccesses(ctx context.Context) ([]entity.Access, error) {
	accesses, err := s.repo.GetAccesses(ctx)
	if err != nil {
		return nil, fmt.Errorf("access.Service.GetAccesses: %v", err)
	}

	return accesses, nil
}
