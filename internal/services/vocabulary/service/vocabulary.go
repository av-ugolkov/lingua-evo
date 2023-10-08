package service

import (
	"context"

	"github.com/google/uuid"
)

type repoDict interface {
	AddWordInDictionary(ctx context.Context, userId, originalWord uuid.UUID, translateWord []uuid.UUID, pronunciation string, examples []uuid.UUID) error
}

type DictionarySvc struct {
	repo repoDict
}

func NewService(repo repoDict) *DictionarySvc {
	return &DictionarySvc{
		repo: repo,
	}
}

func (s *DictionarySvc) AddWordInDictionary(
	ctx context.Context,
	userID uuid.UUID,
	origWordId uuid.UUID,
	tranWordId []uuid.UUID,
	pronunciation string,
	examples []uuid.UUID,
) (uuid.UUID, error) {
	err := s.repo.AddWordInDictionary(ctx, userID, origWordId, tranWordId, pronunciation, examples)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Nil, nil
}
