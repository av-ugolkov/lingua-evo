package dictionary

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	msgerror "github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	entityLanguage "github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
)

type (
	repoDictionary interface {
		AddWords(ctx context.Context, words []DictWord) ([]DictWord, error)
		GetWordsByText(ctx context.Context, words []DictWord) ([]DictWord, error)
		GetWords(ctx context.Context, ids []uuid.UUID) ([]DictWord, error)
		UpdateWord(ctx context.Context, w *DictWord) error
		FindWords(ctx context.Context, w *DictWord) ([]uuid.UUID, error)
		DeleteWordByText(ctx context.Context, word *DictWord) error
		GetRandomWord(ctx context.Context, langCode string) (DictWord, error)
	}

	langSvc interface {
		CheckLanguage(ctx context.Context, langCode string) error
		GetAvailableLanguages(ctx context.Context) ([]entityLanguage.Language, error)
	}
)

//go:generate mockery --inpackage --outpkg dictionary --testonly --name "repoDictionary|langSvc"

type Service struct {
	repo    repoDictionary
	langSvc langSvc
}

func NewService(repo repoDictionary, langSvc langSvc) *Service {
	return &Service{
		repo:    repo,
		langSvc: langSvc,
	}
}

func (s *Service) GetOrAddWords(ctx context.Context, inWords []DictWord) ([]DictWord, error) {
	languages, err := s.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		return nil, msgerror.NewError(fmt.Errorf("dictionary.Service.AddWords - get languages: %v", err), msgerror.ErrInternal)
	}

	dictWords := checkWords(inWords, languages)
	if len(dictWords) == 0 {
		return []DictWord{}, nil
	}

	getWords, err := s.repo.GetWordsByText(ctx, dictWords)
	if err != nil {
		return nil, msgerror.NewError(fmt.Errorf("dictionary.Service.AddWords - get words: %v", err), msgerror.ErrInternal)
	}

	addWords, err := s.repo.AddWords(ctx, dictWords)
	if err != nil {
		return nil, msgerror.NewError(fmt.Errorf("dictionary.Service.AddWords: %v", err), msgerror.ErrInternal)
	}

	words := make([]DictWord, 0, len(getWords)+len(addWords))
	words = append(words, getWords...)
	words = append(words, addWords...)

	return words, nil
}

func (s *Service) GetWordsByID(ctx context.Context, wordIDs []uuid.UUID) ([]DictWord, error) {
	if len(wordIDs) == 0 {
		return []DictWord{}, nil
	}

	words, err := s.repo.GetWords(ctx, wordIDs)
	if err != nil {
		return nil, msgerror.NewError(fmt.Errorf("dictionary.Service.GetWords: %v", err), msgerror.ErrInternal)
	}

	return words, nil
}

func (s *Service) GetWordsByText(ctx context.Context, inWords []DictWord) ([]DictWord, error) {
	languages, err := s.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		return nil, msgerror.NewError(fmt.Errorf("dictionary.Service.AddWords - get languages: %v", err), msgerror.ErrInternal)
	}

	dictWords := checkWords(inWords, languages)
	if len(dictWords) == 0 {
		return nil, msgerror.NewError(fmt.Errorf("dictionary.Service.AddWords - no words"), msgerror.ErrNotFound)
	}

	words, err := s.repo.GetWordsByText(ctx, dictWords)
	if err != nil {
		return nil, msgerror.NewError(fmt.Errorf("dictionary.Service.GetWordByText: %v", err), msgerror.ErrInternal)
	}
	return words, nil
}

func (s *Service) GetRandomWord(ctx context.Context, langCode string) (DictWord, error) {
	err := s.langSvc.CheckLanguage(ctx, langCode)
	if err != nil {
		return DictWord{}, msgerror.NewError(fmt.Errorf("dictionary.Service.GetRandomWord - check language: %v", err), msgerror.ErrInternal)
	}

	word, err := s.repo.GetRandomWord(ctx, langCode)
	if err != nil {
		return DictWord{}, msgerror.NewError(fmt.Errorf("dictionary.Service.GetRandomWord: %v", err), msgerror.ErrInternal)
	}

	return word, nil
}

func (s *Service) DeleteWordByText(ctx context.Context, word DictWord) error {
	err := s.repo.DeleteWordByText(ctx, &word)
	if err != nil {
		return fmt.Errorf("dictionary.Service.DeleteWordByText: %w", err)
	}

	return nil
}

func checkWords(words []DictWord, languages []entityLanguage.Language) []DictWord {
	for i := 0; i < len(words); {
		if !slices.ContainsFunc(languages, func(language entityLanguage.Language) bool {
			return words[i].LangCode == language.Code
		}) {
			slog.Warn("dictionary.checkWords - not validate language")
			words = slices.Delete(words, i, i+1)
			continue
		}

		words[i].Text = strings.TrimSpace(words[i].Text)
		if words[i].Text == runtime.EmptyString {
			slog.Warn("dictionary.checkWords - empty text")
			words = slices.Delete(words, i, i+1)
			continue
		}
		i++
	}

	return words
}
