package example

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type repoExample interface {
	AddExamples(ctx context.Context, ids []uuid.UUID, texts []string, langCode string) ([]uuid.UUID, error)
	GetExampleByValue(ctx context.Context, text string, langCode string) (uuid.UUID, error)
	GetExampleById(ctx context.Context, id uuid.UUID, langCode string) (string, error)
	GetExamples(ctx context.Context, exampleIDs []uuid.UUID) ([]Example, error)
}

type Service struct {
	repo repoExample
}

func NewService(repo repoExample) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddExamples(ctx context.Context, texts []string, langCode string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(texts))
	for i := 0; i < len(texts); i++ {
		texts[i] = strings.TrimSpace(texts[i])
		if texts[i] == "" {
			texts = slices.Delete(texts, i, i+1)
			i--
			continue
		}
		ids = append(ids, uuid.New())
	}
	ids, err := s.repo.AddExamples(ctx, ids, texts, langCode)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *Service) GetExampleById(ctx context.Context, id uuid.UUID, langCode string) (string, error) {
	text, err := s.repo.GetExampleById(ctx, id, langCode)
	if err != nil {
		return "", err
	}

	return text, nil
}

func (s *Service) GetExamples(ctx context.Context, exampleIDs []uuid.UUID) ([]Example, error) {
	if len(exampleIDs) == 0 {
		return []Example{}, nil
	}
	examples, err := s.repo.GetExamples(ctx, exampleIDs)
	if err != nil {
		return nil, fmt.Errorf("example.Service.GetExamples: %w", err)
	}
	return examples, nil
}
