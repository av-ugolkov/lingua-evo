package vocabulary

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/services/vocabulary/model"

	"github.com/google/uuid"
)

type (
	repoVocab interface {
		Add(ctx context.Context, dict Vocabulary, tagIDs []uuid.UUID) error
		Delete(ctx context.Context, dict Vocabulary) error
		GetByName(ctx context.Context, uid uuid.UUID, name string) (Vocabulary, error)
		GetTagsVocabulary(ctx context.Context, vocabularyID uuid.UUID) ([]string, error)
		GetByID(ctx context.Context, dictID uuid.UUID) (Vocabulary, error)
		GetVocabularies(ctx context.Context, userID uuid.UUID) ([]Vocabulary, error)
		Rename(ctx context.Context, id uuid.UUID, newName string) error
	}

	tagSvc interface {
		AddTags(ctx context.Context, tags []string) ([]uuid.UUID, error)
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

func (s *Service) AddVocabulary(ctx context.Context, userID uuid.UUID, data model.VocabularyRq) (Vocabulary, error) {
	vocabulary := Vocabulary{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          data.Name,
		NativeLang:    data.NativeLang,
		TranslateLang: data.TranslateLang,
		Tags:          data.Tags,
	}

	vocabularies, err := s.repoVocab.GetVocabularies(ctx, vocabulary.UserID)
	if err != nil {
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.AddVocabulary - get count vocabularies: %w", err)
	}

	for _, dict := range vocabularies {
		if dict.Name == vocabulary.Name {
			return Vocabulary{}, fmt.Errorf("vocabulary.Service.AddVocabulary - already have vocabulary with same")
		}
	}

	if len(vocabularies) >= 5 {
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.AddVocabulary - %w %v", ErrCountVocabulary, vocabulary.UserID)
	}

	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		tagIDs, err := s.tagSvc.AddTags(ctx, data.Tags)
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
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.AddVocabulary: %w", err)
	}

	vocabulary.NativeLang, err = s.langSvc.GetLangByCode(ctx, vocabulary.NativeLang)
	if err != nil {
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.AddVocabulary - get native lang: %w", err)
	}
	vocabulary.TranslateLang, err = s.langSvc.GetLangByCode(ctx, vocabulary.TranslateLang)
	if err != nil {
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.AddVocabulary - get translate lang: %w", err)
	}

	return vocabulary, nil
}

func (s *Service) DeleteVocabulary(ctx context.Context, userID uuid.UUID, name string) error {
	dict := Vocabulary{
		UserID: userID,
		Name:   name,
	}

	err := s.repoVocab.Delete(ctx, dict)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.DeleteVocabulary: %w", err)
	}
	return nil
}

func (s *Service) GetVocabulary(ctx context.Context, userID uuid.UUID, name string) (Vocabulary, error) {
	vocab, err := s.repoVocab.GetByName(ctx, userID, name)
	if err != nil {
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.GetVocabulary: %w", err)
	}

	vocab.Tags, err = s.repoVocab.GetTagsVocabulary(ctx, vocab.ID)
	if err != nil {
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.GetVocabulary: %w", err)
	}
	return vocab, nil
}

func (s *Service) GetVocabularyByID(ctx context.Context, vocabID uuid.UUID) (Vocabulary, error) {
	vocab, err := s.repoVocab.GetByID(ctx, vocabID)
	if err != nil {
		return Vocabulary{}, fmt.Errorf("vocabulary.Service.GetVocabularyByID: %w", err)
	}

	return vocab, nil
}

func (s *Service) GetVocabularies(ctx context.Context, userID uuid.UUID) ([]Vocabulary, error) {
	vocabularies, err := s.repoVocab.GetVocabularies(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetVocabularies: %w", err)
	}

	return vocabularies, nil
}

func (s *Service) RenameVocabulary(ctx context.Context, id uuid.UUID, newName string) error {
	err := s.repoVocab.Rename(ctx, id, newName)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.RenameVocabulary: %w", err)
	}
	return nil
}
