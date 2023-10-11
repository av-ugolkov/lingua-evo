package service

import (
	"context"

	"lingua-evo/internal/services/vocabulary/dto"
	"lingua-evo/internal/services/vocabulary/entity"
	entityWord "lingua-evo/internal/services/word/entity"

	"github.com/google/uuid"
)

type (
	repoDict interface {
		AddWord(ctx context.Context, vocabulary entity.Vocabulary) error
	}

	wordSvc interface {
		AddWord(ctx context.Context, word *entityWord.Word) (uuid.UUID, error)
	}
)

type VocabularySvc struct {
	repo    repoDict
	wordSvc wordSvc
}

func NewService(repo repoDict, wordSvc wordSvc) *VocabularySvc {
	return &VocabularySvc{
		repo:    repo,
		wordSvc: wordSvc,
	}
}

func (s *VocabularySvc) AddWordInVocabulary(
	ctx context.Context,
	vocab *dto.AddWordRequest,
) (uuid.UUID, error) {
	v := entity.Vocabulary{
		DictionaryId:  uuid.UUID{},
		OriginalWord:  uuid.UUID{},
		TranslateWord: []uuid.UUID{},
		Examples:      []uuid.UUID{},
		Tags:          []uuid.UUID{},
	}

	err := s.repo.AddWord(ctx, v)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Nil, nil
}
