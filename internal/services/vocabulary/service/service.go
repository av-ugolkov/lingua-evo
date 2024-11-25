package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime/access"
	"github.com/av-ugolkov/lingua-evo/tools/math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		GetVocabsWithCountWords(ctx context.Context, uid, owner uuid.UUID, access []uint8) ([]entity.VocabWithUser, error)
		GetWithCountWords(ctx context.Context, vid uuid.UUID) (entity.VocabWithUser, error)
		GetVocabulariesWithMaxWords(ctx context.Context, access []uint8, limit int) ([]entity.VocabWithUser, error)
		GetVocabulariesRecommended(ctx context.Context, uid uuid.UUID, access []uint8, limit uint) ([]entity.VocabWithUser, error)

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

//go:generate mockery --inpackage --outpkg vocabulary --testonly --name "repoVocab|userSvc|exampleSvc|dictSvc|tagSvc|subscribersSvc|eventsSvc"

type Service struct {
	tr             *transactor.Transactor
	repoVocab      repoVocab
	exampleSvc     exampleSvc
	dictSvc        dictSvc
	tagSvc         tagSvc
	subscribersSvc subscribersSvc
	eventsSvc      eventsSvc
}

func NewService(
	tr *transactor.Transactor,
	repoVocab repoVocab,
	exampleSvc exampleSvc,
	dictSvc dictSvc,
	tagSvc tagSvc,
	subscribersSvc subscribersSvc,
	eventsSvc eventsSvc,
) *Service {
	return &Service{
		tr:             tr,
		repoVocab:      repoVocab,
		exampleSvc:     exampleSvc,
		dictSvc:        dictSvc,
		tagSvc:         tagSvc,
		subscribersSvc: subscribersSvc,
		eventsSvc:      eventsSvc,
	}
}

func (s *Service) GetVocabularies(ctx context.Context, uid uuid.UUID, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string, limitWords int) ([]entity.VocabWithUserAndWords, int, error) {
	countItems, err := s.repoVocab.GetVocabulariesCountByAccess(ctx, uid, []access.Type{access.Subscribers, access.Public}, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	if countItems == 0 {
		return []entity.VocabWithUserAndWords{}, 0, nil
	}

	vocabularies, err := s.repoVocab.GetVocabulariesByAccess(ctx, uid, []access.Type{access.Subscribers, access.Public}, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	vocabsWithWords := make([]entity.VocabWithUserAndWords, 0, len(vocabularies))
	for _, v := range vocabularies {
		words, err := s.GetSeveralWords(ctx, v.UserID, v.ID, limitWords)
		if err != nil {
			slog.Error(fmt.Sprintf("vocabulary.Service.GetVocabularies: GetWords: %v", err))
		}

		vocabWords := make([]string, 0, limitWords)
		for _, w := range words[0:math.MinInt(len(words), limitWords)] {
			vocabWords = append(vocabWords, w.Native.Text)
		}

		vocabsWithWords = append(vocabsWithWords, entity.VocabWithUserAndWords{Words: vocabWords, VocabWithUser: v})
	}
	return vocabsWithWords, countItems, nil
}

func (s *Service) GetVocabulary(ctx context.Context, uid, vid uuid.UUID) (entity.Vocab, error) {
	accessStatus, err := s.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return entity.Vocab{},
			msgerr.New(fmt.Errorf("vocabulary.Service.GetVocabulary - %w: %w", entity.ErrAccessDenied, err),
				msgerr.ErrMsgForbidden)
	}

	if accessStatus == access.Forbidden {
		return entity.Vocab{},
			msgerr.New(fmt.Errorf("vocabulary.Service.GetVocabulary - %w", entity.ErrAccessDenied),
				msgerr.ErrMsgForbidden)
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

func (s *Service) GetVocabulariesByUser(ctx context.Context, uid, owner uuid.UUID, access []access.Type) ([]entity.VocabWithUser, error) {
	accessIDs := make([]uint8, len(access))
	for i, v := range access {
		accessIDs[i] = uint8(v)
	}

	vocabs, err := s.repoVocab.GetVocabsWithCountWords(ctx, uid, owner, accessIDs)
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

	if vocab.UserID != uid {
		vocab.Editable, err = s.VocabularyEditable(ctx, uid, vid)
		if err != nil {
			switch {
			case !errors.Is(err, pgx.ErrNoRows):
				return vocab, fmt.Errorf("word.Service.GetWords - check access: %w", err)
			}
		}
	} else {
		vocab.Editable = true
	}

	return vocab, nil
}

func (s *Service) GetRecommendedVocabularies(ctx context.Context, uid uuid.UUID) ([]entity.VocabWithUser, error) {
	if uid == uuid.Nil {
		vocabs, err := s.repoVocab.GetVocabulariesWithMaxWords(ctx, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.GetRecommendedVocabularies: %w", err)
		}
		return vocabs, nil
	}

	vocabs, err := s.repoVocab.GetVocabulariesRecommended(ctx, uid, []uint8{uint8(access.Public), uint8(access.Subscribers)}, 3)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetRecommendedVocabularies: %w", err)
	}

	return vocabs, nil
}
