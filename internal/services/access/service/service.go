package service

import (
	"context"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/access"
	"github.com/av-ugolkov/lingua-evo/internal/services/access/dto"
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

func (s *Service) GetAccessesDTO(ctx context.Context) ([]dto.AccessRs, error) {
	accesses, err := s.repo.GetAccesses(ctx)
	if err != nil {
		return nil, fmt.Errorf("access.Service.GetAccesses: %v", err)
	}

	accessesRs := dto.AccessesToDto(accesses)

	return accessesRs, nil
}
