package service

import (
	"context"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

type (
	repoVocabUser interface {
		GetVocabulariesByUser(ctx context.Context, uid uuid.UUID) ([]entity.Vocab, error)
	}
)

func (s *Service) UserAddVocabulary(ctx context.Context, vocabulary entity.Vocab) (entity.Vocab, error) {
	vocabularies, err := s.repoVocab.GetVocabulariesByUser(ctx, vocabulary.UserID)
	if err != nil {
		return entity.Vocab{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary - get count vocabularies: %w", err)
	}

	for _, dict := range vocabularies {
		if dict.Name == vocabulary.Name {
			return entity.Vocab{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary - already have vocabulary with same")
		}
	}

	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		tagIDs, err := s.tagSvc.AddTags(ctx, vocabulary.Tags)
		if err != nil {
			return fmt.Errorf("add tags: %w", err)
		}

		vocabulary.ID, err = s.repoVocab.AddVocab(ctx, vocabulary, tagIDs)
		if err != nil {
			return fmt.Errorf("add vocabulary: %w", err)
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

func (s *Service) UserGetVocabularies(ctx context.Context, uid uuid.UUID) ([]entity.Vocab, error) {
	vocabularies, err := s.repoVocab.GetVocabulariesByUser(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.UserGetVocabularies: %w", err)
	}

	return vocabularies, nil
}

func (s *Service) UserEditVocabulary(ctx context.Context, vocab entity.Vocab) error {
	err := s.repoVocab.EditVocab(ctx, vocab)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.UserEditVocabulary: %w", err)
	}
	return nil
}
