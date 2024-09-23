package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
)

type (
	repoVocab interface {
		AddVocab(ctx context.Context, vocab entity.Vocab, tagIDs []uuid.UUID) (uuid.UUID, error)
		DeleteVocab(ctx context.Context, vocab entity.Vocab) error
		GetVocab(ctx context.Context, vid uuid.UUID) (entity.Vocab, error)
		GetByName(ctx context.Context, uid uuid.UUID, name string) (entity.Vocab, error)
		GetTagsVocabulary(ctx context.Context, vid uuid.UUID) ([]string, error)
		EditVocab(ctx context.Context, vocab entity.Vocab) error
		GetVocabulariesByAccess(ctx context.Context, uid uuid.UUID, access []access.Type, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabWithUser, error)
		GetVocabulariesCountByAccess(ctx context.Context, uid uuid.UUID, access []access.Type, search, nativeLang, translateLang string) (int, error)
		GetAccess(ctx context.Context, vid uuid.UUID) (uint8, error)
		GetCreatorVocab(ctx context.Context, vid uuid.UUID) (uuid.UUID, error)
		CopyVocab(ctx context.Context, uid, vid uuid.UUID) (uuid.UUID, error)
		GetVocabsWithCountWords(ctx context.Context, uid uuid.UUID, access []uint8) ([]entity.VocabWithUser, error)
		GetWithCountWords(ctx context.Context, vid uuid.UUID) (entity.VocabWithUser, error)
		GetVocabulariesWithMaxWords(ctx context.Context, limit int, access []uint8) ([]entity.VocabWithUser, error)

		repoVocabUser
		repoWord
		repoVocabAccess
	}

	tagSvc interface {
		AddTags(ctx context.Context, tags []entityTag.Tag) ([]uuid.UUID, error)
	}

	subscribersSvc interface {
		Check(ctx context.Context, uid, subID uuid.UUID) (bool, error)
	}
)

type Service struct {
	tr             *transactor.Transactor
	repoVocab      repoVocab
	userSvc        userSvc
	exampleSvc     exampleSvc
	dictSvc        dictSvc
	tagSvc         tagSvc
	subscribersSvc subscribersSvc
}

func NewService(
	tr *transactor.Transactor,
	repoVocab repoVocab,
	userSvc userSvc,
	exampleSvc exampleSvc,
	dictSvc dictSvc,
	tagSvc tagSvc,
	subscribersSvc subscribersSvc,
) *Service {
	return &Service{
		tr:             tr,
		repoVocab:      repoVocab,
		userSvc:        userSvc,
		exampleSvc:     exampleSvc,
		dictSvc:        dictSvc,
		tagSvc:         tagSvc,
		subscribersSvc: subscribersSvc,
	}
}

func (s *Service) GetVocabularies(ctx context.Context, uid uuid.UUID, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabWithUser, int, error) {
	countItems, err := s.repoVocab.GetVocabulariesCountByAccess(ctx, uid, []access.Type{access.Subscribers, access.Public}, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	if countItems == 0 {
		return []entity.VocabWithUser{}, 0, nil
	}

	vocabularies, err := s.repoVocab.GetVocabulariesByAccess(ctx, uid, []access.Type{access.Subscribers, access.Public}, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	return vocabularies, countItems, nil
}

func (s *Service) GetVocabulary(ctx context.Context, uid, vid uuid.UUID) (entity.Vocab, error) {
	accessStatus, err := s.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return entity.Vocab{}, handler.NewError(fmt.Errorf("vocabulary.Service.GetVocabulary - %w: %w", entity.ErrAccessDenied, err),
			http.StatusForbidden, handler.ErrForbidden)
	}

	if accessStatus == access.Forbidden {
		return entity.Vocab{}, handler.NewError(fmt.Errorf("vocabulary.Service.GetVocabulary - %w", entity.ErrAccessDenied),
			http.StatusForbidden, handler.ErrForbidden)
	}

	vocab, err := s.repoVocab.GetVocab(ctx, vid)
	if err != nil {
		return entity.Vocab{}, fmt.Errorf("vocabulary.Service.GetVocabulary: %w", err)
	}

	tags, err := s.repoVocab.GetTagsVocabulary(ctx, vocab.ID)
	if err != nil {
		return entity.Vocab{}, fmt.Errorf("vocabulary.Service.GetVocabulary: %w", err)
	}
	for _, tag := range tags {
		vocab.Tags = append(vocab.Tags, entityTag.Tag{Text: tag})
	}
	return vocab, nil
}

func (s *Service) GetAccessForUser(ctx context.Context, uid, vid uuid.UUID) (access.Status, error) {
	accessID, err := s.repoVocab.GetAccess(ctx, vid)
	if err != nil {
		return access.Forbidden, fmt.Errorf("vocabulary.Service.GetAccessForUser - get access type: %w", err)
	}
	acc := access.Type(accessID)
	if uid == uuid.Nil && acc != access.Public {
		return access.Forbidden, fmt.Errorf("vocabulary.Service.GetAccessForUser - %w", entity.ErrAccessDenied)
	}

	creatodID, err := s.repoVocab.GetCreatorVocab(ctx, vid)
	if err != nil {
		return access.Forbidden, fmt.Errorf("vocabulary.Service.GetAccessForUser - get creator: %w", err)
	}
	if creatodID == uid {
		return access.Edit, nil
	}

	if acc == access.Public {
		return access.Read, nil
	} else {
		editable, err := s.VocabularyEditable(ctx, uid, vid)
		if err != nil {
			return access.Forbidden, fmt.Errorf("vocabulary.Service.GetAccessForUser - get vocabulary access: %w", err)
		}

		isSubscribers, err := s.subscribersSvc.Check(ctx, uid, creatodID)
		if err != nil {
			return access.Forbidden, fmt.Errorf("vocabulary.Service.GetAccessForUser - check subscribers: %w", err)
		}
		if isSubscribers {
			if editable {
				return access.Edit, nil
			}
			return access.Read, nil
		}
	}

	return access.Forbidden, nil
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
	if access.Type(accessID) == access.Public {
		return false, nil
	}

	return true, nil
}

func (s *Service) CopyVocab(ctx context.Context, uid, vid uuid.UUID) error {
	copyVid, err := s.repoVocab.CopyVocab(ctx, uid, vid)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.Copy - copy vocabulary: %w", err)
	}

	err = s.CopyWords(ctx, vid, copyVid)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.Copy - copy words: %w", err)
	}

	return nil
}

func (s *Service) GetVocabulariesByUser(ctx context.Context, uid uuid.UUID, access []access.Type) ([]entity.VocabWithUser, error) {
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

func (s *Service) GetVocabularyInfo(ctx context.Context, uid, vid uuid.UUID) (entity.VocabWithUser, error) {
	vocab, err := s.repoVocab.GetWithCountWords(ctx, vid)
	if err != nil {
		return entity.VocabWithUser{}, fmt.Errorf("vocabulary.Service.GetVocabularyInfo: %w", err)
	}

	return vocab, nil
}

func (s *Service) GetRecommendedVocabularies(ctx context.Context, uid uuid.UUID) ([]entity.VocabWithUser, error) {
	if uid == uuid.Nil {
		vocabs, err := s.repoVocab.GetVocabulariesWithMaxWords(ctx, 3, []uint8{uint8(access.Public), uint8(access.Subscribers)})
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.GetRecommendedVocabularies: %w", err)
		}
		return vocabs, nil
	}

	vocabs, err := s.repoVocab.GetVocabsWithCountWords(ctx, uid, []uint8{1, 2, 3})
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetRecommendedVocabularies: %w", err)
	}

	return vocabs, nil
}
