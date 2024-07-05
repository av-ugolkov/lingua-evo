package access

import (
	"context"
	"fmt"
)

type (
	repoAccess interface {
		GetAccesses(ctx context.Context) ([]Access, error)
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

func (s *Service) GetAccesses(ctx context.Context) ([]Access, error) {
	accesses, err := s.repo.GetAccesses(ctx)
	if err != nil {
		return nil, fmt.Errorf("access.Service.GetAccesses: %w", err)
	}

	return accesses, nil
}
