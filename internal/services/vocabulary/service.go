package vocabulary

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type (
	repoVocabulary interface {
		AddWord(ctx context.Context, vocabulary Vocabulary) error
		DeleteWord(ctx context.Context, vocabulary Vocabulary) error
		GetWords(ctx context.Context, dictID uuid.UUID, limit int) ([]Vocabulary, error)
		UpdateWord(ctx context.Context, vocabulary Vocabulary) error
	}

	exampleSvc interface {
		AddExample(ctx context.Context, text, langCode string) (uuid.UUID, error)
		UpdateExample(ctx context.Context, text, langCode string) (uuid.UUID, error)
	}

	tagSvc interface {
		AddTag(ctx context.Context, text string) (uuid.UUID, error)
		UpdateTag(ctx context.Context, text string) (uuid.UUID, error)
	}

	wordSvc interface {
		AddWord(ctx context.Context, text, langCode, pronunciation string) (uuid.UUID, error)
		UpdateWord(ctx context.Context, text, langCode, pronunciation string) (uuid.UUID, error)
	}
)

type Service struct {
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
) *Service {
	return &Service{
		repo:       repo,
		wordSvc:    wordSvc,
		exampleSvc: exexampleSvc,
		tagSvc:     tagSvc,
	}
}

func (s *Service) AddWord(
	ctx context.Context,
	dictID uuid.UUID,
	nativeWord Word,
	tanslateWords Words,
	examples []string,
	tags []string) (*Vocabulary, error) {
	nativeWordID, err := s.wordSvc.AddWord(ctx, nativeWord.Text, nativeWord.LangCode, nativeWord.Pronunciation)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.AddWord - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(tanslateWords))
	for _, translateWord := range tanslateWords {
		translateWordID, err := s.wordSvc.AddWord(ctx, translateWord.Text, translateWord.LangCode, translateWord.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.AddWord - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}

	exampleIDs := make([]uuid.UUID, 0, len(examples))
	for _, example := range examples {
		exampleID, err := s.exampleSvc.AddExample(ctx, example, nativeWord.LangCode)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.AddWord - add example: %w", err)
		}
		exampleIDs = append(exampleIDs, exampleID)
	}

	tagIDs := make([]uuid.UUID, 0, len(tags))
	for _, tag := range tags {
		tagID, err := s.tagSvc.AddTag(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.AddWord - add tag: %w", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	vocabulary := Vocabulary{
		DictionaryId:   dictID,
		NativeWord:     nativeWordID,
		TranslateWords: translateWordIDs,
		Examples:       exampleIDs,
		Tags:           tagIDs,
	}

	err = s.repo.AddWord(ctx, vocabulary)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.AddWord - add vocabulary: %w", err)
	}

	return &vocabulary, nil
}

func (s *Service) UpdateWord(ctx context.Context,
	dictID uuid.UUID,
	oldWordID uuid.UUID,
	nativeWord Word,
	tanslateWords Words,
	examples []string,
	tags []string) (*VocabularyWord, error) {
	nativeWordID, err := s.wordSvc.UpdateWord(ctx, nativeWord.Text, nativeWord.LangCode, nativeWord.Pronunciation)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(tanslateWords))
	for _, translateWord := range tanslateWords {
		translateWordID, err := s.wordSvc.UpdateWord(ctx, translateWord.Text, translateWord.LangCode, translateWord.Pronunciation)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}

	exampleIDs := make([]uuid.UUID, 0, len(examples))
	for _, example := range examples {
		exampleID, err := s.exampleSvc.UpdateExample(ctx, example, nativeWord.LangCode)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add example: %w", err)
		}
		exampleIDs = append(exampleIDs, exampleID)
	}

	tagIDs := make([]uuid.UUID, 0, len(tags))
	for _, tag := range tags {
		tagID, err := s.tagSvc.UpdateTag(ctx, tag)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add tag: %w", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	vocabulary := Vocabulary{
		DictionaryId:   dictID,
		NativeWord:     nativeWordID,
		TranslateWords: translateWordIDs,
		Examples:       exampleIDs,
		Tags:           tagIDs,
	}

	if oldWordID != nativeWordID {
		err = s.repo.DeleteWord(ctx, Vocabulary{DictionaryId: dictID, NativeWord: oldWordID})
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - delete old word: %w", err)
		}
		err = s.repo.AddWord(ctx, vocabulary)
		if err != nil {
			return nil, fmt.Errorf("vocabulary.Service.UpdateWord - add new word: %w", err)
		}
		return &VocabularyWord{
			NativeWord:     nativeWord.Text,
			TranslateWords: tanslateWords.GetValues(),
			Examples:       examples,
			Tags:           tags,
		}, nil
	}
	err = s.repo.UpdateWord(ctx, vocabulary)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.UpdateWord - update vocabulary: %w", err)
	}

	return &VocabularyWord{
		NativeWord:     nativeWord.Text,
		TranslateWords: tanslateWords.GetValues(),
		Examples:       examples,
		Tags:           tags,
	}, nil
}

func (s *Service) DeleteWord(ctx context.Context, dictID, nativeWordID uuid.UUID) error {
	vocabulary := Vocabulary{
		DictionaryId: dictID,
		NativeWord:   nativeWordID,
	}

	err := s.repo.DeleteWord(ctx, vocabulary)
	if err != nil {
		return fmt.Errorf("vocabulary.Service.DeleteWord - delete word: %w", err)
	}
	return nil
}

func (s *Service) GetWords(ctx context.Context, dictID uuid.UUID, limit int) ([]VocabularyWord, error) {
	vocabularies, err := s.repo.GetWords(ctx, dictID, limit)
	if err != nil {
		return nil, fmt.Errorf("vocabulary.Service.GetWords - get words: %w", err)
	}

	for vocabulary := range vocabularies {
		slog.Info(fmt.Sprintf("%+v", vocabulary))
	}

	words := []VocabularyWord{}
	return words, nil
}
