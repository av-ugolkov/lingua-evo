package service

import (
	"context"
	"errors"
	"fmt"

	"lingua-evo/internal/services/lingua/language/entity"
)

type langRepo interface {
	GetAvailableLanguages(ctx context.Context) ([]*entity.Language, error)
	GetLanguage(ctx context.Context, lang string) (*entity.Language, error)
}

type LanguageSvc struct {
	repo langRepo
}

func NewService(repo langRepo) *LanguageSvc {
	return &LanguageSvc{
		repo: repo,
	}
}

func (s *LanguageSvc) GetLanguage(ctx context.Context, lang string) (*entity.Language, error) {
	language, err := s.repo.GetLanguage(ctx, lang)
	if err != nil {
		return nil, fmt.Errorf("service.lingua.GetLanguage: %v", err)
	}

	return language, nil
}

func (s *LanguageSvc) GetAvailableLanguages(ctx context.Context) ([]*entity.Language, error) {
	languages, err := s.repo.GetAvailableLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.lingua.GetAvailableLanguages: %v", err)
	}

	return languages, nil
}

func (s *LanguageSvc) CheckLanguage(ctx context.Context, langCode string) error {
	if len(langCode) == 0 {
		return errors.New("language.service.LanguageSvc.CheckLanguage - code language is empty")
	}

	if _, err := s.repo.GetLanguage(ctx, langCode); err != nil {
		return fmt.Errorf("language.service.LanguageSvc.CheckLanguage - not found language: %s", langCode)
	}
	return nil
}
