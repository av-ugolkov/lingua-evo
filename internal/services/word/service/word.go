package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"lingua-evo/internal/services/word/entity"
)

type repoWord interface {
	AddWord(ctx context.Context, w *entity.Word) (uuid.UUID, error)
	GetWord(ctx context.Context, text, langCode string) (*entity.Word, error)
	EditWord(ctx context.Context, w *entity.Word) error
	FindWord(ctx context.Context, w *entity.Word) (uuid.UUID, error)
	RemoveWord(ctx context.Context, w *entity.Word) error
	GetRandomWord(ctx context.Context, lang string) (*entity.Word, error)
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

func (s *WordSvc) AddWord(ctx context.Context, word *entity.Word) (uuid.UUID, error) {
	repoWord, err := s.repo.GetWord(ctx, word.Text, word.LanguageCode)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("word.service.WordSvc.AddWord - get word: %w", err)
	} else if repoWord != nil {
		return repoWord.ID, nil
	}
	wordID, err := s.repo.AddWord(ctx, word)
	if err != nil {
		return uuid.Nil, err
	}

	return wordID, nil
}

func (s *WordSvc) EditWord(ctx context.Context, w *entity.Word) error {
	return nil
}

func (s *WordSvc) FindWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	return nil, nil
}

func (s *WordSvc) FindWords(ctx context.Context, w string) (*entity.Word, error) {
	return nil, nil
}

func (s *WordSvc) RemoveWord(ctx context.Context, w *entity.Word) error {
	return nil
}

func (s *WordSvc) GetRandomWord(ctx context.Context, lang string) (*entity.Word, error) {
	return s.repo.GetRandomWord(ctx, lang)
}

func (s *WordSvc) SharedWord(ctx context.Context, w *entity.Word) (*entity.Word, error) {
	return nil, nil
}
