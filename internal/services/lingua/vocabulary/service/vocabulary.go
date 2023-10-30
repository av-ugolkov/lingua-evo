package service

import (
	"context"
	"fmt"

	"lingua-evo/internal/services/lingua/vocabulary/dto"
	"lingua-evo/internal/services/lingua/vocabulary/entity"
	dtoWord "lingua-evo/internal/services/lingua/word/dto"

	"github.com/google/uuid"
)

type (
	repoVocabulary interface {
		AddWord(ctx context.Context, vocabulary entity.Vocabulary) error
		DeleteWord(ctx context.Context, vocabulary entity.Vocabulary) (int64, error)
	}

	exampleSvc interface {
		AddExample(ctx context.Context, text, langCode string) (uuid.UUID, error)
	}

	tagSvc interface {
		AddTag(ctx context.Context, text string) (uuid.UUID, error)
	}

	wordSvc interface {
		AddWord(ctx context.Context, word *dtoWord.AddWordRq) (uuid.UUID, error)
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

func (s *VocabularySvc) AddWordInVocabulary(ctx context.Context, v *dto.AddWordRq) error {
	word := dtoWord.AddWordRq{
		Text:          v.NativeWord.Text,
		Pronunciation: v.NativeWord.Pronunciation,
		LanguageCode:  v.NativeWord.LangCode,
	}
	nativeWordID, err := s.wordSvc.AddWord(ctx, &word)
	if err != nil {
		return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(v.TanslateWords))
	for _, translateWord := range v.TanslateWords {
		word = dtoWord.AddWordRq{
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

	exampleIDs := make([]uuid.UUID, 0, len(v.Examples))
	for _, example := range v.Examples {
		exampleID, err := s.exampleSvc.AddExample(ctx, example, v.NativeWord.LangCode)
		if err != nil {
			return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add example: %w", err)
		}
		exampleIDs = append(exampleIDs, exampleID)
	}

	tagIDs := make([]uuid.UUID, 0, len(v.Tags))
	for _, tag := range v.Tags {
		tagID, err := s.tagSvc.AddTag(ctx, tag)
		if err != nil {
			return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add tag: %w", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	vocabulary := entity.Vocabulary{
		DictionaryId:   v.DictionaryID,
		NativeWord:     nativeWordID,
		TranslateWords: translateWordIDs,
		Examples:       exampleIDs,
		Tags:           tagIDs,
	}

	err = s.repo.AddWord(ctx, vocabulary)
	if err != nil {
		return fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add vocabulary: %w", err)
	}

	return nil
}

func (s *VocabularySvc) DeleteWordFromVocabulary(ctx context.Context, v *dto.RemoveWordRq) error {
	vocabulary := entity.Vocabulary{
		DictionaryId: v.DictionaryID,
		NativeWord:   v.NativeWordID,
	}

	_, err := s.repo.DeleteWord(ctx, vocabulary)
	if err != nil {
		return fmt.Errorf("vocabulary.service.VocabularySvc.DeleteWordFromVocabulary - delete word: %w", err)
	}
	return nil
}
