package access

import (
	"context"
	"fmt"

	msgerror "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
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
		return nil, msgerror.NewError(fmt.Errorf("access.Service.GetAccesses: %v", err), msgerror.ErrInternal)
	}

	return accesses, nil
}
