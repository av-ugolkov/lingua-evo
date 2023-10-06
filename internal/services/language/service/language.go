package service

import (
	"context"
	"errors"
	"fmt"

	"lingua-evo/internal/services/language/entity"
)

type langRepo interface {
	GetLanguages(context.Context) ([]*entity.Language, error)
	GetNameLanguage(ctx context.Context, lang string) (string, error)
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

func (s *LanguageSvc) CheckLanguage(ctx context.Context, codeLang string) error {
	if len(codeLang) == 0 {
		return errors.New("language.service.LanguageSvc.CheckLanguage - code language is empty")
	}

	if _, err := s.repo.GetNameLanguage(ctx, codeLang); err != nil {
		return fmt.Errorf("language.service.LanguageSvc.CheckLanguage - not found language: %s", codeLang)
	}
	return nil
}
