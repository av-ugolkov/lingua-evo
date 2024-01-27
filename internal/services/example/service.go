package example

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type repoExample interface {
	AddExample(ctx context.Context, id uuid.UUID, text, langCode string) error
	GetExample(ctx context.Context, text string, langCode string) (uuid.UUID, error)
	GetExampleById(ctx context.Context, id uuid.UUID, langCode string) (string, error)
}

type Service struct {
	repo repoExample
}

func NewService(repo repoExample) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddExample(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	id := uuid.New()
	err := s.repo.AddExample(ctx, id, text, langCode)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *Service) GetExample(ctx context.Context, id uuid.UUID, langCode string) (string, error) {
	text, err := s.repo.GetExampleById(ctx, id, langCode)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (s *Service) UpdateExample(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	id, err := s.repo.GetExample(ctx, text, langCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("example.Service.UpdateExample: %w", err)
	}
	if id != uuid.Nil {
		return id, nil
	}
	id = uuid.New()
	err = s.repo.AddExample(ctx, id, text, langCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("example.Service.UpdateExample: %w", err)
	}
	return id, nil
}
