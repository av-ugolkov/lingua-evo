package service

import (
	"context"
	"fmt"

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
		Get(ctx context.Context, vid uuid.UUID) (entity.Vocabulary, error)
		GetByName(ctx context.Context, uid uuid.UUID, name string) (entity.Vocabulary, error)
		GetTagsVocabulary(ctx context.Context, vid uuid.UUID) ([]string, error)
		GetByID(ctx context.Context, vid uuid.UUID) (entity.Vocabulary, error)
		GetVocabulariesByUser(ctx context.Context, uid uuid.UUID) ([]entity.Vocabulary, error)
		Edit(ctx context.Context, vocab entity.Vocabulary) error
		GetVocabulariesByAccess(ctx context.Context, uid uuid.UUID, access []int, page, itemsPerPage, typeOrder int, search string) ([]entity.Vocabulary, error)
		GetVocabulariesCountByAccess(ctx context.Context, uid uuid.UUID, access []int, search string) (int, error)
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

func (s *Service) GetVocabularies(ctx context.Context, uid uuid.UUID, page, itemsPerPage, typeOrder int, search string) ([]entity.Vocabulary, int, error) {
	countItems, err := s.repoVocab.GetVocabulariesCountByAccess(ctx, uid, []int{AccessSubscribers, AccessPublic}, search)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	if countItems == 0 {
		return []entity.Vocabulary{}, 0, nil
	}

	vocabularies, err := s.repoVocab.GetVocabulariesByAccess(ctx, uid, []int{AccessSubscribers, AccessPublic}, page, itemsPerPage, typeOrder, search)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	return vocabularies, countItems, nil
}
