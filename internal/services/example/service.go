package example

import (
	"context"
	"fmt"
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

func (s *Service) AddExamples(ctx context.Context, examples []Example, langCode string) ([]uuid.UUID, error) {
	if len(examples) == 0 {
		return []uuid.UUID{}, nil
	}

	ids := make([]uuid.UUID, 0, len(examples))
	texts := make([]string, 0, len(examples))
	for i := 0; i < len(examples); i++ {
		text := strings.TrimSpace(examples[i].Text)
		if text == "" {
			continue
		}
		texts = append(texts, text)
		ids = append(ids, examples[i].ID)
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
