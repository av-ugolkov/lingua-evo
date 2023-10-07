package service

import (
	"context"

	"github.com/google/uuid"

	"lingua-evo/internal/services/word/entity"
)

type repoWord interface {
	AddWord(ctx context.Context, w *entity.Word) (uuid.UUID, error)
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
