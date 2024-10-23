package language

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	"github.com/av-ugolkov/lingua-evo/runtime"
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
		return runtime.EmptyString, msgerr.New(fmt.Errorf("language.Service.GetLanguage: %v", err), msgerr.ErrMsgInternal)
	}

	return language, nil
}

func (s *Service) GetAvailableLanguages(ctx context.Context) ([]Language, error) {
	languages, err := s.repo.GetAvailableLanguages(ctx)
	if err != nil {
		return nil, msgerr.New(fmt.Errorf("language.Service.GetAvailableLanguages: %v", err), msgerr.ErrMsgInternal)
	}

	return languages, nil
}

func (s *Service) CheckLanguage(ctx context.Context, langCode string) error {
	if len(langCode) == 0 {
		return nil
	}

	if _, err := s.repo.GetLanguage(ctx, langCode); err != nil {
		return msgerr.New(fmt.Errorf("language.Service.CheckLanguage: %w", err), msgerr.ErrMsgInternal)
	}
	return nil
}
