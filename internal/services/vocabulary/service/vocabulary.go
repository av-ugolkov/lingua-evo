package service

import (
	"context"
	"fmt"

	"lingua-evo/internal/services/vocabulary/dto"
	"lingua-evo/internal/services/vocabulary/entity"
	entityWord "lingua-evo/internal/services/word/entity"

	"github.com/google/uuid"
)

type (
	repoDict interface {
		AddWord(ctx context.Context, vocabulary entity.Vocabulary) error
	}

	exampleSvc interface {
		AddExample(ctx context.Context, text, langCode string) (uuid.UUID, error)
	}

	tagSvc interface {
		AddTag(ctx context.Context, text string) (uuid.UUID, error)
	}

	wordSvc interface {
		AddWord(ctx context.Context, word *entityWord.Word) (uuid.UUID, error)
	}
)

type VocabularySvc struct {
	repo       repoDict
	wordSvc    wordSvc
	exampleSvc exampleSvc
	tagSvc     tagSvc
}

func NewService(
	repo repoDict,
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

func (s *VocabularySvc) AddWordInVocabulary(ctx context.Context, vocab *dto.AddWordRq) error {
	word := entityWord.Word{
		ID:            uuid.New(),
		Text:          vocab.NativeWord.Text,
		Pronunciation: vocab.NativeWord.Pronunciation,
		LanguageCode:  vocab.NativeWord.LangCode,
	}
	nativeWordID, err := s.wordSvc.AddWord(ctx, &word)
	if err != nil {
		return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(vocab.TanslateWords))
	for _, translateWord := range vocab.TanslateWords {
		word = entityWord.Word{
			ID:            uuid.New(),
			Text:          translateWord.Text,
			Pronunciation: translateWord.Pronunciation,
			LanguageCode:  translateWord.LangCode,
		}
		translateWordID, err := s.wordSvc.AddWord(ctx, &word)
		if err != nil {
			return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}

	exampleIDs := make([]uuid.UUID, 0, len(vocab.Examples))
	for _, example := range vocab.Examples {
		exampleID, err := s.exampleSvc.AddExample(ctx, example, vocab.NativeWord.LangCode)
		if err != nil {
			return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add example: %w", err)
		}
		exampleIDs = append(exampleIDs, exampleID)
	}

	tagIDs := make([]uuid.UUID, 0, len(vocab.Tags))
	for _, tag := range vocab.Tags {
		tagID, err := s.tagSvc.AddTag(ctx, tag)
		if err != nil {
			return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add tag: %w", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	v := entity.Vocabulary{
		DictionaryId:  vocab.DictionaryID,
		NativeWord:    nativeWordID,
		TranslateWord: translateWordIDs,
		Examples:      exampleIDs,
		Tags:          tagIDs,
	}

	err = s.repo.AddWord(ctx, v)
	if err != nil {
		return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add vocabulary: %w", err)
	}

	return nil
}
