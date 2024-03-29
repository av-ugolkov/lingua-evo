package tag

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type (
	repoTag interface {
		AddTag(ctx context.Context, id uuid.UUID, text string) (uuid.UUID, error)
		FindTag(ctx context.Context, text string) ([]*Tag, error)
		GetTag(ctx context.Context, text string) (uuid.UUID, error)
		GetTags(ctx context.Context, tagIDs []uuid.UUID) ([]Tag, error)
		GetAllTags(ctx context.Context) ([]Tag, error)
	}
)

type Service struct {
	repo repoTag
}

func NewService(repo repoTag) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddTag(ctx context.Context, text string) (uuid.UUID, error) {
	id, err := s.repo.AddTag(ctx, uuid.New(), text)
	if err != nil {
		return uuid.Nil, fmt.Errorf("tag.Service.AddTag: %w", err)
	}
	return id, nil
}

func (s *Service) FindTag(ctx context.Context, text string) ([]*Tag, error) {
	tags, err := s.repo.FindTag(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("tag.Service.FindTag: %w", err)
	}
	return tags, nil
}

func (s *Service) GetAllTag(ctx context.Context) ([]Tag, error) {
	tags, err := s.repo.GetAllTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("tag.Service.GetAllTag: %w", err)
	}
	return tags, nil
}

func (s *Service) GetTags(ctx context.Context, tagIDs []uuid.UUID) ([]Tag, error) {
	if len(tagIDs) == 0 {
		return []Tag{}, nil
	}
	tags, err := s.repo.GetTags(ctx, tagIDs)
	if err != nil {
		return nil, fmt.Errorf("tag.Service.GetTags: %w", err)
	}
	return tags, nil
}

func (s *Service) UpdateTag(ctx context.Context, text string) (uuid.UUID, error) {
	id, err := s.repo.GetTag(ctx, text)
	if err != nil {
		return uuid.Nil, fmt.Errorf("tag.Service.UpdateTag: %w", err)
	}
	if id != uuid.Nil {
		return id, nil
	}

	id, err = s.AddTag(ctx, text)
	if err != nil {
		return uuid.Nil, fmt.Errorf("tag.Service.UpdateTag: %w", err)
	}
	return id, nil
}
