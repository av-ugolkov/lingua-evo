package service

import (
	"context"
	"fmt"
	"slices"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entityAccess "github.com/av-ugolkov/lingua-evo/internal/services/access"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

const (
	AccessPrivate     = 0
	AccessSubscribers = 1
	AccessPublic      = 2
)

type (
	repoVocab interface {
		Add(ctx context.Context, vocab entity.Vocabulary, tagIDs []uuid.UUID) error
		Delete(ctx context.Context, vocab entity.Vocabulary) error
		Get(ctx context.Context, vocabID uuid.UUID) (entity.Vocabulary, error)
		GetByName(ctx context.Context, uid uuid.UUID, name string) (entity.Vocabulary, error)
		GetTagsVocabulary(ctx context.Context, vocabID uuid.UUID) ([]string, error)
		GetByID(ctx context.Context, vocabID uuid.UUID) (entity.Vocabulary, error)
		GetVocabulariesByUser(ctx context.Context, userID uuid.UUID) ([]entity.Vocabulary, error)
		Edit(ctx context.Context, vocab entity.Vocabulary) error
		GetVocabulariesByAccess(ctx context.Context, access []int) ([]entity.Vocabulary, error)
	}

	repoAccess interface {
		GetAccesses(ctx context.Context) ([]entityAccess.Access, error)
	}

	tagSvc interface {
		AddTags(ctx context.Context, tags []entityTag.Tag) ([]uuid.UUID, error)
	}

	langSvc interface {
		GetLangByCode(ctx context.Context, code string) (string, error)
	}
)

type Service struct {
	tr        *transactor.Transactor
	repoVocab repoVocab
	langSvc   langSvc
	tagSvc    tagSvc
}

func NewService(tr *transactor.Transactor, repoVocab repoVocab, langSvc langSvc, tagSvc tagSvc) *Service {
	return &Service{
		tr:        tr,
		repoVocab: repoVocab,
		langSvc:   langSvc,
		tagSvc:    tagSvc,
	}
}

func (s *Service) GetVocabularies(ctx context.Context, userID uuid.UUID) ([]entity.Vocabulary, error) {
	vocabularies, err := s.repoVocab.GetVocabulariesByAccess(ctx, []int{AccessSubscribers, AccessPublic})
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	if userID != uuid.Nil {
		userVocabs, err := s.repoVocab.GetVocabulariesByUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
		}

		for _, vocab := range userVocabs {
			if !slices.ContainsFunc(vocabularies, func(v entity.Vocabulary) bool {
				return vocab.ID == v.ID
			}) {
				vocabularies = append(vocabularies, vocab)
			}
		}
	}

	return vocabularies, nil
}
