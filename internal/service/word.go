package service

import (
	"context"

	"github.com/google/uuid"

	"lingua-evo/internal/delivery/repository"
)

type WordsService struct {
	WordDB repository.WordDB
}

func NewWordsService(db repository.WordDB) *WordsService {
	return &WordsService{
		WordDB: db,
	}
}

func (s *WordsService) SendWord(ctx context.Context, w *repository.Word) (uuid.UUID, error) {
	wordId, err := s.WordDB.AddWord(ctx, w)
	if err != nil {
		return uuid.UUID{}, err
	}
	return wordId, nil
}

func (s *WordsService) EditWord(ctx context.Context, w *repository.Word) error {
	return nil
}

func (s *WordsService) FindWord(ctx context.Context, w string) (*repository.Word, error) {
	return nil, nil
}

func (s *WordsService) RemoveWord(ctx context.Context, w *repository.Word) error {
	return nil
}

func (s *WordsService) PickRandomWord(ctx context.Context, w *repository.Word) (*repository.Word, error) {
	return nil, nil
}

func (s *WordsService) SharedWord(ctx context.Context, w *repository.Word) (*repository.Word, error) {
	return nil, nil
}
