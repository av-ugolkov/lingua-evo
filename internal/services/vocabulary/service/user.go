package service

import (
	"context"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime/access"

	"github.com/google/uuid"
)

type (
	repoVocabUser interface {
		GetVocabulariesByUser(ctx context.Context, uid uuid.UUID) ([]entity.VocabWithUser, error)
		GetVocabulariesCountByUser(ctx context.Context, uid uuid.UUID, access []access.Type, search, nativeLang, translateLang string) (int, error)
		GetSortedVocabulariesByUser(ctx context.Context, uid uuid.UUID, access []access.Type, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabWithUser, error)
	}
)

func (s *Service) UserAddVocabulary(ctx context.Context, vocabulary entity.Vocab) (entity.Vocab, error) {
	vocabularies, err := s.repoVocab.GetVocabulariesByUser(ctx, vocabulary.UserID)
	if err != nil {
		return entity.Vocab{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary - get count vocabularies: %w", err)
	}

	for _, vocab := range vocabularies {
		if vocab.Name == vocabulary.Name {
			return entity.Vocab{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary - already have vocabulary with same")
		}
	}

	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		id, err := s.repoVocab.AddVocab(ctx, vocabulary)
		if err != nil {
			return fmt.Errorf("add vocabulary: %w", err)
		}

		vocabulary, err = s.repoVocab.GetVocab(ctx, id)
		if err != nil {
			return fmt.Errorf("get vocabulary: %w", err)
		}
		return nil
	})

	if err != nil {
		return entity.Vocab{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary: %w", err)
	}

	return vocabulary, nil
}

func (s *Service) UserDeleteVocabulary(ctx context.Context, userID uuid.UUID, name string) error {
	dict := entity.Vocab{
		UserID: userID,
		Name:   name,
	}

	err := s.repoVocab.DeleteVocab(ctx, dict)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.UserDeleteVocabulary: %w", err)
	}
	return nil
}

func (s *Service) UserGetVocabularies(ctx context.Context, uid uuid.UUID, page, itemsPerPage, typeSort, order int, search, nativeLang, translateLang string) ([]entity.VocabWithUser, int, error) {
	countItems, err := s.repoVocab.GetVocabulariesCountByUser(ctx, uid, []access.Type{access.Subscribers, access.Public}, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.UserGetVocabularies: %w", err)
	}

	if countItems == 0 {
		return []entity.VocabWithUser{}, 0, nil
	}

	vocabs, err := s.repoVocab.GetSortedVocabulariesByUser(ctx, uid, []access.Type{access.Subscribers, access.Public}, page, itemsPerPage, typeSort, order, search, nativeLang, translateLang)
	if err != nil {
		return nil, 0, fmt.Errorf("vocabulary.Service.UserGetVocabularies: %w", err)
	}

	return vocabs, countItems, nil
}

func (s *Service) UserEditVocabulary(ctx context.Context, vocab entity.Vocab) error {
	err := s.repoVocab.EditVocab(ctx, vocab)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.UserEditVocabulary: %w", err)
	}
	return nil
}
