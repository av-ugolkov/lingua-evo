package service

import (
	"context"
	"fmt"

	entityTag "github.com/av-ugolkov/lingua-evo/internal/services/tag"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

func (s *Service) UserAddVocabulary(ctx context.Context, vocabulary entity.Vocabulary) (entity.Vocabulary, error) {
	vocabularies, err := s.repoVocab.GetVocabulariesByUser(ctx, vocabulary.UserID)
	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary - get count vocabularies: %w", err)
	}

	for _, dict := range vocabularies {
		if dict.Name == vocabulary.Name {
			return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary - already have vocabulary with same")
		}
	}

	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		tagIDs, err := s.tagSvc.AddTags(ctx, vocabulary.Tags)
		if err != nil {
			return fmt.Errorf("add tags: %w", err)
		}

		err = s.repoVocab.Add(ctx, vocabulary, tagIDs)
		if err != nil {
			return fmt.Errorf("add vocabulary: %w", err)
		}

		return nil
	})

	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.UserAddVocabulary: %w", err)
	}

	return vocabulary, nil
}

func (s *Service) UserDeleteVocabulary(ctx context.Context, userID uuid.UUID, name string) error {
	dict := entity.Vocabulary{
		UserID: userID,
		Name:   name,
	}

	err := s.repoVocab.Delete(ctx, dict)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.UserDeleteVocabulary: %w", err)
	}
	return nil
}

func (s *Service) UserGetVocabulary(ctx context.Context, vocabID uuid.UUID) (entity.Vocabulary, error) {
	vocab, err := s.repoVocab.Get(ctx, vocabID)
	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.UserGetVocabulary: %w", err)
	}

	tags, err := s.repoVocab.GetTagsVocabulary(ctx, vocab.ID)
	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.UserGetVocabulary: %w", err)
	}
	for _, tag := range tags {
		vocab.Tags = append(vocab.Tags, entityTag.Tag{Text: tag})
	}
	return vocab, nil
}

func (s *Service) UserGetVocabularyByID(ctx context.Context, vocabID uuid.UUID) (entity.Vocabulary, error) {
	vocab, err := s.repoVocab.GetByID(ctx, vocabID)
	if err != nil {
		return entity.Vocabulary{}, fmt.Errorf("vocabulary.Service.UserGetVocabularyByID: %w", err)
	}

	return vocab, nil
}

func (s *Service) UserGetVocabularies(ctx context.Context, userID uuid.UUID) ([]entity.Vocabulary, error) {
	vocabularies, err := s.repoVocab.GetVocabulariesByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.UserGetVocabularies: %w", err)
	}

	return vocabularies, nil
}

func (s *Service) UserEditVocabulary(ctx context.Context, vocab entity.Vocabulary) error {
	err := s.repoVocab.Edit(ctx, vocab)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.UserEditVocabulary: %w", err)
	}
	return nil
}
