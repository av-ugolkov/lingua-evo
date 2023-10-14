package service

import (
	"context"
	"fmt"

	"lingua-evo/internal/services/tag/entity"

	"github.com/google/uuid"
)

type (
	repoTag interface {
		AddTag(ctx context.Context, id uuid.UUID, text string) (uuid.UUID, error)
		FindTag(ctx context.Context, text string) ([]*entity.Tag, error)
		GetAllTags(ctx context.Context) ([]*entity.Tag, error)
	}
)

type TagSvc struct {
	repo repoTag
}

func NewService(repo repoTag) *TagSvc {
	return &TagSvc{
		repo: repo,
	}
}

func (s *TagSvc) AddTag(ctx context.Context, text string) (uuid.UUID, error) {
	id, err := s.repo.AddTag(ctx, uuid.New(), text)
	if err != nil {
		return uuid.Nil, fmt.Errorf("tag.service.TagSvc.AddTag: %w", err)
	}
	return id, nil
}

func (s *TagSvc) FindTag(ctx context.Context, text string) ([]*entity.Tag, error) {
	tags, err := s.repo.FindTag(ctx, text)
	if err != nil {
		return nil, fmt.Errorf("tag.service.TagSvc.FindTag: %w", err)
	}
	return tags, nil
}

func (s *TagSvc) GetAllTag(ctx context.Context) ([]*entity.Tag, error) {
	tags, err := s.repo.GetAllTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("tag.service.TagSvc.GetAllTag: %w", err)
	}
	return tags, nil
}
