package dictionary

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strings"

	entityLanguage "github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	repoDictionary interface {
		AddWords(ctx context.Context, words []DictWord) ([]uuid.UUID, error)
		GetWordIDByText(ctx context.Context, w *DictWord) (uuid.UUID, error)
		GetWords(ctx context.Context, ids []uuid.UUID) ([]DictWord, error)
		UpdateWord(ctx context.Context, w *DictWord) error
		FindWords(ctx context.Context, w *DictWord) ([]uuid.UUID, error)
		DeleteWordByID(ctx context.Context, id uuid.UUID) (int64, error)
		DeleteWordByText(ctx context.Context, text, langCode string) (int64, error)
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

func (s *Service) AddWords(ctx context.Context, words []DictWord) ([]uuid.UUID, error) {
	if len(words) == 0 {
		return nil, fmt.Errorf("dictionary.Service.AddWords - empty word list")
	}

	languages, err := s.langSvc.GetAvailableLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.AddWords - get languages: %v", err)
	}
	for i := 0; i < len(words); {
		word := words[i]
		if !slices.ContainsFunc(languages, func(language *entityLanguage.Language) bool {
			return word.LangCode == language.Code
		}) {
			slog.Warn(fmt.Sprintf("dictionary.Service.AddWords - check language: %v", err))
			words = slices.Delete(words, i, i+1)
			continue
		}

		word.Text = strings.TrimSpace(word.Text)
		if word.Text == runtime.EmptyString {
			slog.Warn("dictionary.Service.AddWords - empty text")
			words = slices.Delete(words, i, i+1)
		}
		i++
	}
	wordIDs, err := s.repo.AddWords(ctx, words)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.AddWords: %v", err)
	}

	return wordIDs, nil
}

func (s *Service) GetWordByText(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	text = strings.TrimSpace(text)
	if text == runtime.EmptyString {
		return uuid.Nil, fmt.Errorf("dictionary.Service.GetWordByText - empty text")
	}

	if err := s.langSvc.CheckLanguage(ctx, langCode); err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.GetWordByText - check language: %v", err)
	}

	word := DictWord{
		Text:     text,
		LangCode: langCode,
	}

	wordID, err := s.repo.GetWordIDByText(ctx, &word)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.GetWordByText: %v", err)
	}
	return wordID, nil
}

func (s *Service) GetWords(ctx context.Context, wordIDs []uuid.UUID) ([]DictWord, error) {
	if len(wordIDs) == 0 {
		return []DictWord{}, nil
	}
	words, err := s.repo.GetWords(ctx, wordIDs)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetWords: %w", err)
	}
	return words, nil
}

func (s *Service) FindWords(ctx context.Context, text, langCode string) ([]uuid.UUID, error) {
	word := DictWord{
		Text:     text,
		LangCode: langCode,
	}

	wordIDs, err := s.repo.FindWords(ctx, &word)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.FindWord: %w", err)
	}

	return wordIDs, nil
}

func (s *Service) DeleteWord(ctx context.Context, text, langCode string) error {
	i, err := s.repo.DeleteWordByText(ctx, text, langCode)
	if err != nil {
		return fmt.Errorf("dictionary.Service.DeleteWord: %w", err)
	}

	slog.Debug(fmt.Sprintf("deleted %d rows", i))

	return nil
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

func (s *Service) UpdateWord(ctx context.Context, word DictWord) (uuid.UUID, error) {
	wordID, err := s.repo.GetWordIDByText(ctx, &word)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("dictionary.Service.UpdateWord - get word: %w", err)
	} else if wordID == uuid.Nil {
		//TODO check correct work
		wordIDs, err := s.repo.AddWords(ctx, []DictWord{word})
		if err != nil {
			return uuid.Nil, fmt.Errorf("dictionary.Service.UpdateWord - add word: %w", err)
		}

		return wordIDs[0], nil
	}

	word.ID = wordID

	err = s.repo.UpdateWord(ctx, &word)
	if err != nil {
		return uuid.Nil, err
	}

	return wordID, nil
}
