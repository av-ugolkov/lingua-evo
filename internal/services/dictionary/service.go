package dictionary

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/av-ugolkov/lingua-evo/internal/services/dictionary/model"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type (
	repoDictionary interface {
		AddWords(ctx context.Context, words []Word) ([]uuid.UUID, error)
		GetWordByText(ctx context.Context, w *Word) (uuid.UUID, error)
		GetWords(ctx context.Context, ids []uuid.UUID) ([]Word, error)
		UpdateWord(ctx context.Context, w *Word) error
		FindWords(ctx context.Context, w *Word) ([]uuid.UUID, error)
		DeleteWord(ctx context.Context, w *Word) (int64, error)
		GetRandomWord(ctx context.Context, w *Word) (*Word, error)
	}

	langSvc interface {
		CheckLanguage(ctx context.Context, langCode string) error
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

func (s *Service) AddWord(ctx context.Context, data model.WordRq) (uuid.UUID, error) {
	data.Text = strings.TrimSpace(data.Text)
	if data.Text == runtime.EmptyString {
		return uuid.Nil, fmt.Errorf("dictionary.Service.AddWord - empty text")
	}

	err := s.langSvc.CheckLanguage(ctx, data.LangCode)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.AddWord - check language: %v", err)
	}

	wordIDs, err := s.repo.AddWords(ctx, []Word{
		{
			ID:            uuid.New(),
			Text:          data.Text,
			LangCode:      data.LangCode,
			Pronunciation: data.Pronunciation,
		},
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.AddWord: %v", err)
	}

	return wordIDs[0], nil
}

func (s *Service) GetWordByText(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	text = strings.TrimSpace(text)
	if text == runtime.EmptyString {
		return uuid.Nil, fmt.Errorf("dictionary.Service.GetWordByText - empty text")
	}

	if err := s.langSvc.CheckLanguage(ctx, langCode); err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.GetWordByText - check language: %v", err)
	}

	word := Word{
		Text:     text,
		LangCode: langCode,
	}

	wordID, err := s.repo.GetWordByText(ctx, &word)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.GetWordByText: %v", err)
	}
	return wordID, nil
}

func (s *Service) GetWords(ctx context.Context, wordIDs []uuid.UUID) ([]Word, error) {
	if len(wordIDs) == 0 {
		return []Word{}, nil
	}
	words, err := s.repo.GetWords(ctx, wordIDs)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetWords: %w", err)
	}
	return words, nil
}

func (s *Service) EditWord(ctx context.Context, text, langCode string) error {
	word := Word{
		Text:     text,
		LangCode: langCode,
	}
	fmt.Println(word)

	return nil
}

func (s *Service) FindWords(ctx context.Context, text, langCode string) ([]uuid.UUID, error) {
	word := Word{
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
	word := Word{
		Text:     text,
		LangCode: langCode,
	}

	i, err := s.repo.DeleteWord(ctx, &word)
	if err != nil {
		return fmt.Errorf("dictionary.Service.DeleteWord: %w", err)
	}

	slog.Debug(fmt.Sprintf("deleted %d rows", i))

	return nil
}

func (s *Service) GetRandomWord(ctx context.Context, langCode string) (*Word, error) {
	if err := s.langSvc.CheckLanguage(ctx, langCode); err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetRandomWord - check language: %v", err)
	}

	word := &Word{
		LangCode: langCode,
	}

	word, err := s.repo.GetRandomWord(ctx, word)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetRandomWord: %w", err)
	}

	return word, nil
}

func (s *Service) UpdateWord(ctx context.Context, text, langCode, pronunciation string) (uuid.UUID, error) {
	word := &Word{
		Text:          text,
		LangCode:      langCode,
		Pronunciation: pronunciation,
	}

	wordID, err := s.repo.GetWordByText(ctx, word)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("dictionary.Service.UpdateWord - get word: %w", err)
	} else if wordID == uuid.Nil {
		word.ID = uuid.New()
		wordIDs, err := s.repo.AddWords(ctx, []Word{*word})
		if err != nil {
			return uuid.Nil, fmt.Errorf("dictionary.Service.UpdateWord - add word: %w", err)
		}

		return wordIDs[0], nil
	}

	word.ID = wordID

	err = s.repo.UpdateWord(ctx, word)
	if err != nil {
		return uuid.Nil, err
	}

	return wordID, nil
}

func (s *Service) SharedWord(ctx context.Context, w *Word) (*Word, error) {
	slog.Error("dictionary.Service.SharedWord - not implemented")
	return nil, nil
}
