package service

import (
	"context"

	"github.com/google/uuid"
)

type repoExample interface {
	AddExample(ctx context.Context, wordId uuid.UUID, example string) (uuid.UUID, error)
}

type ExampleSvc struct {
	repo repoExample
}

func NewService(repo repoExample) *ExampleSvc {
	return &ExampleSvc{
		repo: repo,
	}
}

func (s *ExampleSvc) AddExample(ctx context.Context, wordId uuid.UUID, example string) (uuid.UUID, error) {
	exampleId, err := s.repo.AddExample(ctx, wordId, example)
	if err != nil {
		return uuid.Nil, err
	}

	return exampleId, nil
}
