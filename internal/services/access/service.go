package access

import (
	"context"
	"fmt"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
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
		return nil, handler.NewError(fmt.Errorf("access.Service.GetAccesses: %v", err), handler.ErrInternal)
	}

	return accesses, nil
}
