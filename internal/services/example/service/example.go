package service

import (
	"context"

	"github.com/google/uuid"
)

type repoExample interface {
	AddExample(ctx context.Context, id uuid.UUID, text, langCode string) error
	GetExample(ctx context.Context, id uuid.UUID, langCode string) (string, error)
}

type ExampleSvc struct {
	repo repoExample
}

func NewService(repo repoExample) *ExampleSvc {
	return &ExampleSvc{
		repo: repo,
	}
}

func (s *ExampleSvc) AddExample(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	id := uuid.New()
	err := s.repo.AddExample(ctx, id, text, langCode)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *ExampleSvc) GetExample(ctx context.Context, id uuid.UUID, langCode string) (string, error) {
	text, err := s.repo.GetExample(ctx, id, langCode)
	if err != nil {
		return "", err
	}

	return text, nil
}
