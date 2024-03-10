package word

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type repoWord interface {
	AddWord(ctx context.Context, w *Word) (uuid.UUID, error)
	GetWordByText(ctx context.Context, w *Word) (uuid.UUID, error)
	GetWords(ctx context.Context, ids []uuid.UUID) ([]Word, error)
	UpdateWord(ctx context.Context, w *Word) error
	FindWords(ctx context.Context, w *Word) ([]uuid.UUID, error)
	DeleteWord(ctx context.Context, w *Word) (int64, error)
	GetRandomWord(ctx context.Context, w *Word) (*Word, error)
	SharedWord(ctx context.Context, w *Word) (*Word, error)
}

//go:generate mockery --inpackage --outpkg word --testonly --name "repoWord"

type Service struct {
	repo repoWord
}

func NewService(repo repoWord) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddWord(ctx context.Context, id uuid.UUID, text, langCode, pronunciation string) (uuid.UUID, error) {
	word := &Word{
		ID:            id,
		Text:          text,
		LanguageCode:  langCode,
		Pronunciation: pronunciation,
	}

	id, err := s.repo.AddWord(ctx, word)
	if err != nil {
		return uuid.Nil, fmt.Errorf("word.Service.AddWord: %w", err)
	}

	return id, nil
}

func (s *Service) GetWordByValue(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	word := Word{
		Text:         text,
		LanguageCode: langCode,
	}

	wordID, err := s.repo.GetWordByText(ctx, &word)
	if err != nil {
		return uuid.Nil, fmt.Errorf("word.Service.GetWordByValue: %w", err)
	}
	return wordID, nil
}

func (s *Service) GetWords(ctx context.Context, wordIDs []uuid.UUID) ([]Word, error) {
	if len(wordIDs) == 0 {
		return []Word{}, nil
	}
	words, err := s.repo.GetWords(ctx, wordIDs)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords: %w", err)
	}
	return words, nil
}

func (s *Service) EditWord(ctx context.Context, text, langCode string) error {
	word := Word{
		Text:         text,
		LanguageCode: langCode,
	}
	fmt.Println(word)

	return nil
}

func (s *Service) FindWords(ctx context.Context, text, langCode string) ([]uuid.UUID, error) {
	word := Word{
		Text:         text,
		LanguageCode: langCode,
	}

	wordIDs, err := s.repo.FindWords(ctx, &word)
	if err != nil {
		return nil, fmt.Errorf("word.Service.FindWord: %w", err)
	}

	return wordIDs, nil
}

func (s *Service) DeleteWord(ctx context.Context, text, langCode string) error {
	word := Word{
		Text:         text,
		LanguageCode: langCode,
	}

	i, err := s.repo.DeleteWord(ctx, &word)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord: %w", err)
	}

	slog.Debug(fmt.Sprintf("deleted %d rows", i))

	return nil
}

func (s *Service) GetRandomWord(ctx context.Context, langCode string) (*Word, error) {
	word := &Word{
		LanguageCode: langCode,
	}

	return s.repo.GetRandomWord(ctx, word)
}

func (s *Service) UpdateWord(ctx context.Context, text, langCode, pronunciation string) (uuid.UUID, error) {
	word := &Word{
		Text:          text,
		LanguageCode:  langCode,
		Pronunciation: pronunciation,
	}

	wordID, err := s.repo.GetWordByText(ctx, word)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("word.Service.UpdateWord - get word: %w", err)
	} else if wordID == uuid.Nil {
		word.ID = uuid.New()
		wordID, err = s.repo.AddWord(ctx, word)
		if err != nil {
			return uuid.Nil, fmt.Errorf("word.Service.UpdateWord - add word: %w", err)
		}

		return wordID, nil
	}

	word.ID = wordID

	err = s.repo.UpdateWord(ctx, word)
	if err != nil {
		return uuid.Nil, err
	}

	return wordID, nil
}

func (s *Service) SharedWord(ctx context.Context, w *Word) (*Word, error) {
	slog.Error("word.Service.SharedWord - not implemented")
	return nil, nil
}
