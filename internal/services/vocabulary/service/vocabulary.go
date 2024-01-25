package service

import (
	"context"
	"fmt"

	entity "lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

type (
	repoVocabulary interface {
		AddWord(ctx context.Context, vocabulary entity.Vocabulary) error
		DeleteWord(ctx context.Context, vocabulary entity.Vocabulary) (int64, error)
		GetWords(ctx context.Context, dictID uuid.UUID) ([]entity.Vocabulary, error)
	}

	exampleSvc interface {
		AddExample(ctx context.Context, text, langCode string) (uuid.UUID, error)
	}

	tagSvc interface {
		AddTag(ctx context.Context, text string) (uuid.UUID, error)
	}

	wordSvc interface {
		AddWord(ctx context.Context, text, langCode, pronunciation string) (uuid.UUID, error)
	}
)

type VocabularySvc struct {
	repo       repoVocabulary
	wordSvc    wordSvc
	exampleSvc exampleSvc
	tagSvc     tagSvc
}

func NewService(
	repo repoVocabulary,
	wordSvc wordSvc,
	exexampleSvc exampleSvc,
	tagSvc tagSvc,
) *VocabularySvc {
	return &VocabularySvc{
		repo:       repo,
		wordSvc:    wordSvc,
		exampleSvc: exexampleSvc,
		tagSvc:     tagSvc,
	}
}

func (s *VocabularySvc) AddWordInVocabulary(
	ctx context.Context,
	dictID uuid.UUID,
	nativeWord entity.Word,
	tanslateWords []entity.Word,
	examples []string,
	tags []string) (*entity.Vocabulary, error) {
	nativeWordID, err := s.wordSvc.AddWord(ctx, nativeWord.Text, nativeWord.LangCode, nativeWord.Pronunciation)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(tanslateWords))
	for _, translateWord := range tanslateWords {
		translateWordID, err := s.wordSvc.AddWord(ctx, translateWord.Text, translateWord.LangCode, translateWord.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}

	exampleIDs := make([]uuid.UUID, 0, len(examples))
	for _, example := range examples {
		exampleID, err := s.exampleSvc.AddExample(ctx, example, nativeWord.LangCode)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add example: %w", err)
		}
		exampleIDs = append(exampleIDs, exampleID)
	}

	tagIDs := make([]uuid.UUID, 0, len(tags))
	for _, tag := range tags {
		tagID, err := s.tagSvc.AddTag(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add tag: %w", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	vocabulary := entity.Vocabulary{
		DictionaryId:   dictID,
		NativeWord:     nativeWordID,
		TranslateWords: translateWordIDs,
		Examples:       exampleIDs,
		Tags:           tagIDs,
	}

	err = s.repo.AddWord(ctx, vocabulary)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add vocabulary: %w", err)
	}

	return &vocabulary, nil
}

func (s *VocabularySvc) DeleteWordFromVocabulary(ctx context.Context, dictID, nativeWordID uuid.UUID) error {
	vocabulary := entity.Vocabulary{
		DictionaryId: dictID,
		NativeWord:   nativeWordID,
	}

	_, err := s.repo.DeleteWord(ctx, vocabulary)
	if err != nil {
		return fmt.Errorf("vocabulary.service.VocabularySvc.DeleteWordFromVocabulary - delete word: %w", err)
	}
	return nil
}

func (s *VocabularySvc) GetWords(ctx context.Context, dictID uuid.UUID) ([]entity.Vocabulary, error) {
	words, err := s.repo.GetWords(ctx, dictID)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.service.VocabularySvc.GetWords - get words: %w", err)
	}
	return words, nil
}
