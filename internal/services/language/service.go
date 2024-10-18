package language

import (
	"context"
	"fmt"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
)

type langRepo interface {
	GetAvailableLanguages(ctx context.Context) ([]Language, error)
	GetLanguage(ctx context.Context, lang string) (string, error)
}

type Service struct {
	repo langRepo
}

func NewService(repo langRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetLangByCode(ctx context.Context, lang string) (string, error) {
	language, err := s.repo.GetLanguage(ctx, lang)
	if err != nil {
		return "", handler.NewError(fmt.Errorf("language.Service.GetLanguage: %v", err), handler.ErrInternal)
	}

	return language, nil
}

func (s *Service) GetAvailableLanguages(ctx context.Context) ([]Language, error) {
	languages, err := s.repo.GetAvailableLanguages(ctx)
	if err != nil {
		return nil, handler.NewError(fmt.Errorf("language.Service.GetAvailableLanguages: %v", err), handler.ErrInternal)
	}

	return languages, nil
}

func (s *Service) CheckLanguage(ctx context.Context, langCode string) error {
	if len(langCode) == 0 {
		return nil
	}

	if _, err := s.repo.GetLanguage(ctx, langCode); err != nil {
		return handler.NewError(fmt.Errorf("language.Service.CheckLanguage: %w", err), handler.ErrInternal)
	}
	return nil
}
