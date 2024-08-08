package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

type (
	repoVocab interface {
		AddVocab(ctx context.Context, vocab entity.Vocabulary, tagIDs []uuid.UUID) (uuid.UUID, error)
		DeleteVocab(ctx context.Context, vocab entity.Vocabulary) error
		GetVocab(ctx context.Context, vid uuid.UUID) (entity.Vocabulary, error)
		GetByName(ctx context.Context, uid uuid.UUID, name string) (entity.Vocabulary, error)
		GetTagsVocabulary(ctx context.Context, vid uuid.UUID) ([]string, error)
		GetVocabulariesByUser(ctx context.Context, uid uuid.UUID) ([]entity.Vocabulary, error)
		EditVocab(ctx context.Context, vocab entity.Vocabulary) error
		GetVocabulariesByAccess(ctx context.Context, uid uuid.UUID, access []uint8, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabularyWithUser, error)
		GetVocabulariesCountByAccess(ctx context.Context, uid uuid.UUID, access []uint8, search, nativeLang, translateLang string) (int, error)
		GetAccess(ctx context.Context, vid uuid.UUID) (uint8, error)
		GetCreatorVocab(ctx context.Context, vid uuid.UUID) (uuid.UUID, error)
		CopyVocab(ctx context.Context, uid, vid uuid.UUID) (uuid.UUID, error)
		GetVocabsWithCountWords(ctx context.Context, uid uuid.UUID, access []uint8) ([]entity.VocabularyWithUser, error)
		GetWithCountWords(ctx context.Context, vid uuid.UUID) (entity.VocabularyWithUser, error)
	}

	tagSvc interface {
		AddTags(ctx context.Context, tags []entityTag.Tag) ([]uuid.UUID, error)
	}

	subscribersSvc interface {
		Check(ctx context.Context, uid, subID uuid.UUID) (bool, error)
	}

	vocabAccessSvc interface {
		VocabularyEditable(ctx context.Context, uid, vid uuid.UUID) (bool, error)
	}
)

type Service struct {
	tr             *transactor.Transactor
	repoVocab      repoVocab
	tagSvc         tagSvc
	subscribersSvc subscribersSvc
	vocabAccessSvc vocabAccessSvc
}

func NewService(
	tr *transactor.Transactor,
	repoVocab repoVocab,
	tagSvc tagSvc,
	subscribersSvc subscribersSvc,
	vocabAccessSvc vocabAccessSvc,
) *Service {
	return &Service{
		tr:             tr,
		repoVocab:      repoVocab,
		tagSvc:         tagSvc,
		subscribersSvc: subscribersSvc,
		vocabAccessSvc: vocabAccessSvc,
	}
}

func (s *Service) GetVocabularies(ctx context.Context, uid uuid.UUID, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabularyWithUser, int, error) {
	countItems, err := s.repoVocab.GetVocabulariesCountByAccess(ctx, uid, []uint8{uint8(entity.AccessSubscribers), uint8(entity.AccessPublic)}, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	if countItems == 0 {
		return []entity.VocabularyWithUser{}, 0, nil
	}

	vocabularies, err := s.repoVocab.GetVocabulariesByAccess(ctx, uid, []uint8{uint8(entity.AccessSubscribers), uint8(entity.AccessPublic)}, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	return vocabularies, countItems, nil
}

func (s *Service) GetVocabulary(ctx context.Context, uid, vid uuid.UUID) (entity.Vocabulary, error) {
	err := s.CheckAccess(ctx, uid, vid)
	if err != nil {
		return entity.Vocabulary{}, handler.NewError(fmt.Errorf("vocabulary.Service.GetVocabulary - %w: %w", entity.ErrAccessDenied, err),
			http.StatusForbidden, handler.ErrForbidden)
	}

	vocab, err := s.repoVocab.GetVocab(ctx, vid)
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

func (s *Service) CheckAccess(ctx context.Context, userID, vocabID uuid.UUID) error {
	accessID, err := s.repoVocab.GetAccess(ctx, vocabID)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.checkAccess - get access type: %w", err)
	}
	access := entity.AccessVocab(accessID)
	if userID == uuid.Nil && access != entity.AccessPublic {
		return fmt.Errorf("vocabulary.Service.checkAccess - %w", entity.ErrAccessDenied)
	}
	if access == entity.AccessPublic {
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

	_, err = s.vocabAccessSvc.VocabularyEditable(ctx, userID, vocabID)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.checkAccess - get vocabulary access: %w", err)
	}

	return nil
}

func (s *Service) CanEdit(ctx context.Context, uid, vid uuid.UUID) (bool, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vid)
	if err != nil {
		return false, fmt.Errorf("vocabulary.Service.CanEdit - get vocabulary: %w", err)
	}

	if vocab.UserID == uid {
		return true, nil
	}

	accessID, err := s.repoVocab.GetAccess(ctx, vid)
	if err != nil {
		return false, fmt.Errorf("vocabulary.Service.CanEdit - get access type: %w", err)
	}
	if entity.AccessVocab(accessID) == entity.AccessPublic {
		return false, nil
	}

	return true, nil
}

func (s *Service) CopyVocab(ctx context.Context, uid, vid uuid.UUID) (uuid.UUID, error) {
	copyVid, err := s.repoVocab.CopyVocab(ctx, uid, vid)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.Service.Copy - copy vocabulary: %w", err)
	}

	return copyVid, nil
}

func (s *Service) GetVocabulariesByUser(ctx context.Context, uid uuid.UUID, access []entity.AccessVocab) ([]entity.VocabularyWithUser, error) {
	accessIDs := make([]uint8, len(access))
	for i, v := range access {
		accessIDs[i] = uint8(v)
	}

	vocabs, err := s.repoVocab.GetVocabsWithCountWords(ctx, uid, accessIDs)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetVocabulariesByUser: %w", err)
	}

	return vocabs, nil
}

func (s *Service) GetVocabularyInfo(ctx context.Context, uid, vid uuid.UUID) (entity.VocabularyWithUser, error) {
	vocab, err := s.repoVocab.GetWithCountWords(ctx, vid)
	if err != nil {
		return entity.VocabularyWithUser{}, fmt.Errorf("vocabulary.Service.GetVocabularyInfo: %w", err)
	}

	return vocab, nil
}
