package dictionary

import (
	"context"
	"fmt"
	entityLanguage "github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"log/slog"
	"slices"
	"strings"

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
		GetAvailableLanguages(ctx context.Context) ([]*entityLanguage.Language, error)
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
		return nil, fmt.Errorf("dictionary.Service.AddWords - get languages: %v", err)
	}

	dictWords := checkWords(inWords, languages)
	if len(dictWords) == 0 {
		return []DictWord{}, nil
	}

	getWords, err := s.repo.GetWordsByText(ctx, dictWords)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.AddWords - get words: %v", err)
	}

	addWords, err := s.repo.AddWords(ctx, dictWords)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.AddWords: %v", err)
	}

	words := make([]DictWord, 0, len(getWords)+len(addWords))
	for _, w := range getWords {
		words = append(words, w)
	}
	for _, w := range addWords {
		words = append(words, w)
	}

	return words, nil
}

func (s *Service) GetWordsByID(ctx context.Context, wordIDs []uuid.UUID) ([]DictWord, error) {
	if len(wordIDs) == 0 {
		return []DictWord{}, nil
	}

	words, err := s.repo.GetWords(ctx, wordIDs)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetWords: %w", err)
	}

	return words, nil
}

func (s *Service) GetWordsByText(ctx context.Context, inWords []DictWord) ([]DictWord, error) {
	languages, err := s.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.AddWords - get languages: %v", err)
	}

	dictWords := checkWords(inWords, languages)
	if len(dictWords) == 0 {
		return []DictWord{}, nil
	}

	words, err := s.repo.GetWordsByText(ctx, dictWords)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetWordByText: %v", err)
	}
	return words, nil
}

func (s *Service) GetRandomWord(ctx context.Context, langCode string) (DictWord, error) {
	if err := s.langSvc.CheckLanguage(ctx, langCode); err != nil {
		return DictWord{}, fmt.Errorf("dictionary.Service.GetRandomWord - check language: %v", err)
	}

	word, err := s.repo.GetRandomWord(ctx, langCode)
	if err != nil {
		return DictWord{}, fmt.Errorf("dictionary.Service.GetRandomWord: %w", err)
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

func checkWords(words []DictWord, languages []*entityLanguage.Language) []DictWord {
	for i := 0; i < len(words); {
		if !slices.ContainsFunc(languages, func(language *entityLanguage.Language) bool {
			return words[i].LangCode == language.Code
		}) {
			slog.Warn(fmt.Sprintf("dictionary.checkWords - not validate language"))
			words = slices.Delete(words, i, i+1)
			continue
		}

		words[i].Text = strings.TrimSpace(words[i].Text)
		if words[i].Text == runtime.EmptyString {
			slog.Warn(fmt.Sprintf("dictionary.checkWords - empty text"))
			words = slices.Delete(words, i, i+1)
			continue
		}
		i++
	}

	return words
}
