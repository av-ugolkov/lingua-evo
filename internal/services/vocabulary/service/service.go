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
	AccessPrivate     uint8 = 0
	AccessSubscribers uint8 = 1
	AccessPublic      uint8 = 2
)

type (
	repoVocab interface {
		Add(ctx context.Context, vocab entity.Vocabulary, tagIDs []uuid.UUID) error
		Delete(ctx context.Context, vocab entity.Vocabulary) error
		Get(ctx context.Context, vid uuid.UUID) (entity.Vocabulary, error)
		GetByName(ctx context.Context, uid uuid.UUID, name string) (entity.Vocabulary, error)
		GetTagsVocabulary(ctx context.Context, vid uuid.UUID) ([]string, error)
		GetVocabulariesByUser(ctx context.Context, uid uuid.UUID) ([]entity.Vocabulary, error)
		Edit(ctx context.Context, vocab entity.Vocabulary) error
		GetVocabulariesByAccess(ctx context.Context, uid uuid.UUID, access []uint8, page, itemsPerPage, typeOrder int, search, nativeLang, translateLang string) ([]entity.VocabularyWithUser, error)
		GetVocabulariesCountByAccess(ctx context.Context, uid uuid.UUID, access []uint8, search, nativeLang, translateLang string) (int, error)
		GetAccess(ctx context.Context, vid uuid.UUID) (uint8, error)
		GetCreatorVocab(ctx context.Context, vid uuid.UUID) (uuid.UUID, error)
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

	subscribersSvc interface {
		Check(ctx context.Context, uid, subID uuid.UUID) (bool, error)
	}

	vocabAccessSvc interface {
		GetVocabularyAccess(ctx context.Context, uid, vid uuid.UUID) (bool, error)
	}
)

type Service struct {
	tr             *transactor.Transactor
	repoVocab      repoVocab
	langSvc        langSvc
	tagSvc         tagSvc
	subscribersSvc subscribersSvc
	vocabAccessSvc vocabAccessSvc
}

func NewService(
	tr *transactor.Transactor,
	repoVocab repoVocab,
	langSvc langSvc,
	tagSvc tagSvc,
	subscribersSvc subscribersSvc,
	vocabAccessSvc vocabAccessSvc) *Service {
	return &Service{
		tr:             tr,
		repoVocab:      repoVocab,
		langSvc:        langSvc,
		tagSvc:         tagSvc,
		subscribersSvc: subscribersSvc,
		vocabAccessSvc: vocabAccessSvc,
	}
}

func (s *Service) GetVocabularies(ctx context.Context, uid uuid.UUID, page, itemsPerPage, typeOrder int, search, nativeLang, translateLang string) ([]entity.VocabularyWithUser, int, error) {
	countItems, err := s.repoVocab.GetVocabulariesCountByAccess(ctx, uid, []uint8{AccessSubscribers, AccessPublic}, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	if countItems == 0 {
		return []entity.VocabularyWithUser{}, 0, nil
	}

	vocabularies, err := s.repoVocab.GetVocabulariesByAccess(ctx, uid, []uint8{AccessSubscribers, AccessPublic}, page, itemsPerPage, typeOrder, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	return vocabularies, countItems, nil
}

func (s *Service) GetVocabulary(ctx context.Context, uid, vocabID uuid.UUID) (entity.Vocabulary, error) {
	err := s.checkAccess(ctx, uid, vocabID)
	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.GetVocabulary - %w: %w", entity.ErrAccessDenied, err)
	}

	vocab, err := s.repoVocab.Get(ctx, vocabID)
	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.GetVocabulary: %w", err)
	}

	tags, err := s.repoVocab.GetTagsVocabulary(ctx, vocab.ID)
	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.GetVocabulary: %w", err)
	}
	for _, tag := range tags {
		vocab.Tags = append(vocab.Tags, entityTag.Tag{Text: tag})
	}
	return vocab, nil
}

func (s *Service) checkAccess(ctx context.Context, userID, vocabID uuid.UUID) error {
	accessID, err := s.repoVocab.GetAccess(ctx, vocabID)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.checkAccess - get access type: %w", err)
	}
	if accessID == AccessPublic {
		return nil
	}

	creatodID, err := s.repoVocab.GetCreatorVocab(ctx, vocabID)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.checkAccess - get creator: %w", err)
	}

	if creatodID == userID {
		return nil
	}

	isSubscribers, err := s.subscribersSvc.Check(ctx, creatodID, userID)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.checkAccess - check subscribers: %w", err)
	}
	if isSubscribers {
		return nil
	}

	_, err = s.vocabAccessSvc.GetVocabularyAccess(ctx, userID, vocabID)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.checkAccess - get vocabulary access: %w", err)
	}

	return nil
}
