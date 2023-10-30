package service

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lingua-evo/internal/services/lingua/word/dto"
	"lingua-evo/internal/services/lingua/word/entity"
)

type repoWord interface {
	AddWord(ctx context.Context, w *entity.Word) (uuid.UUID, error)
	GetWord(ctx context.Context, w *entity.Word) (uuid.UUID, error)
	EditWord(ctx context.Context, w *entity.Word) (int64, error)
	FindWords(ctx context.Context, w *entity.Word) ([]uuid.UUID, error)
	DeleteWord(ctx context.Context, w *entity.Word) (int64, error)
	GetRandomWord(ctx context.Context, w *entity.Word) (*entity.Word, error)
	SharedWord(ctx context.Context, w *entity.Word) (*entity.Word, error)
}

type WordSvc struct {
	repo repoWord
}

func NewService(repo repoWord) *WordSvc {
	return &WordSvc{
		repo: repo,
	}
}

func (s *WordSvc) AddWord(ctx context.Context, w *dto.AddWordRq) (uuid.UUID, error) {
	word := &entity.Word{
		ID:            uuid.New(),
		Text:          w.Text,
		LanguageCode:  w.LanguageCode,
		Pronunciation: w.Pronunciation,
	}

	wordID, err := s.repo.GetWord(ctx, word)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("word.service.WordSvc.AddWord - get word: %w", err)
	} else if wordID != uuid.Nil {
		return wordID, nil
	}
	wordID, err = s.repo.AddWord(ctx, word)
	if err != nil {
		return uuid.Nil, err
	}

	return wordID, nil
}

func (s *WordSvc) GetWord(ctx context.Context, w *dto.GetWordRq) (uuid.UUID, error) {
	word := entity.Word{
		Text:         w.Text,
		LanguageCode: w.LanguageCode,
	}

	wordID, err := s.repo.GetWord(ctx, &word)
	if err != nil {
		return uuid.Nil, fmt.Errorf("word.service.WordSvc.GetWord: %w", err)
	}
	return wordID, nil
}

func (s *WordSvc) EditWord(ctx context.Context, w *dto.GetWordRq) error {
	return nil
}

func (s *WordSvc) FindWords(ctx context.Context, w *dto.GetWordRq) ([]uuid.UUID, error) {
	word := entity.Word{
		Text:         w.Text,
		LanguageCode: w.LanguageCode,
	}

	wordIDs, err := s.repo.FindWords(ctx, &word)
	if err != nil {
		return nil, fmt.Errorf("word.service.WordSvc.FindWord: %w", err)
	}

	return wordIDs, nil
}

func (s *WordSvc) DeleteWord(ctx context.Context, w *dto.GetWordRq) error {
	word := entity.Word{
		Text:         w.Text,
		LanguageCode: w.LanguageCode,
	}

	i, err := s.repo.DeleteWord(ctx, &word)
	if err != nil {
		return fmt.Errorf("word.service.WordSvc.DeleteWord: %w", err)
	}

	slog.Debug(fmt.Sprintf("deleted %d rows", i))

	return nil
}

func (s *WordSvc) GetRandomWord(ctx context.Context, w *dto.RandomWordRq) (*entity.Word, error) {
	word := &entity.Word{
		LanguageCode: w.LanguageCode,
	}

	return s.repo.GetRandomWord(ctx, word)
}

func (s *WordSvc) SharedWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	return nil, nil
}