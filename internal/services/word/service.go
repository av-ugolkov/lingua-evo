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
	GetWord(ctx context.Context, w *Word) (uuid.UUID, error)
	UpdateWord(ctx context.Context, w *Word) error
	FindWords(ctx context.Context, w *Word) ([]uuid.UUID, error)
	DeleteWord(ctx context.Context, w *Word) (int64, error)
	GetRandomWord(ctx context.Context, w *Word) (*Word, error)
	SharedWord(ctx context.Context, w *Word) (*Word, error)
}

type Service struct {
	repo repoWord
}

func NewService(repo repoWord) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddWord(ctx context.Context, text, langCode, pronunciation string) (uuid.UUID, error) {
	word := &Word{
		ID:            uuid.New(),
		Text:          text,
		LanguageCode:  langCode,
		Pronunciation: pronunciation,
	}

	wordID, err := s.repo.GetWord(ctx, word)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("word.Service.AddWord - get word: %w", err)
	} else if wordID != uuid.Nil {
		return wordID, nil
	}
	wordID, err = s.repo.AddWord(ctx, word)
	if err != nil {
		return uuid.Nil, err
	}

	return wordID, nil
}

func (s *Service) GetWord(ctx context.Context, text, langCode string) (uuid.UUID, error) {
	word := Word{
		Text:         text,
		LanguageCode: langCode,
	}

	wordID, err := s.repo.GetWord(ctx, &word)
	if err != nil {
		return uuid.Nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}
	return wordID, nil
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

	wordID, err := s.repo.GetWord(ctx, word)

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
	slog.Error("not implemented")
	return nil, nil
}
