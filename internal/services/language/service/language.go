package service

import (
	"context"
	"fmt"

	"lingua-evo/internal/services/language/entity"
)

type langRepo interface {
	GetLanguages(context.Context) ([]*entity.Language, error)
}

type LanguageSvc struct {
	repo langRepo
}

func NewService(repo langRepo) *LanguageSvc {
	return &LanguageSvc{
		repo: repo,
	}
}

func (s *LanguageSvc) GetLanguages(ctx context.Context) ([]*entity.Language, error) {
	languages, err := s.repo.GetLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.lingua.GetLanguages: %v", err)
	}

	return languages, nil
}
