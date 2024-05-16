package word

import (
	"context"
	"errors"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entityVocab "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

type (
	repoWord interface {
		GetWord(ctx context.Context, wordID uuid.UUID) (VocabWordData, error)
		AddWord(ctx context.Context, word VocabWord) error
		DeleteWord(ctx context.Context, word VocabWord) error
		GetRandomVocabulary(ctx context.Context, vocabID uuid.UUID, limit int) ([]VocabWordData, error)
		GetVocabularyWords(ctx context.Context, vocabID uuid.UUID) ([]VocabWordData, error)
		UpdateWord(ctx context.Context, word VocabWord) error
	}

	vocabSvc interface {
		GetVocabularyByID(ctx context.Context, vocabID uuid.UUID) (entityVocab.Vocabulary, error)
	}

	exampleSvc interface {
		AddExamples(ctx context.Context, examples []entityExample.Example, langCode string) ([]uuid.UUID, error)
		GetExamples(ctx context.Context, exampleIDs []uuid.UUID) ([]entityExample.Example, error)
	}

	dictSvc interface {
		GetOrAddWords(ctx context.Context, words []entityDict.DictWord) ([]entityDict.DictWord, error)
		GetWordsByID(ctx context.Context, wordIDs []uuid.UUID) ([]entityDict.DictWord, error)
		GetWordsByText(ctx context.Context, words []entityDict.DictWord) ([]entityDict.DictWord, error)
	}
)

type Service struct {
	tr         *transactor.Transactor
	repo       repoWord
	vocabSvc   vocabSvc
	dictSvc    dictSvc
	exampleSvc exampleSvc
}

func NewService(
	tr *transactor.Transactor,
	repo repoWord,
	vocabSvc vocabSvc,
	dictSvc dictSvc,
	exampleSvc exampleSvc,
) *Service {
	return &Service{
		tr:         tr,
		repo:       repo,
		vocabSvc:   vocabSvc,
		dictSvc:    dictSvc,
		exampleSvc: exampleSvc,
	}
}

func (s *Service) AddWord(ctx context.Context, vocabWordData VocabWordData) (VocabWord, error) {
	vocab, err := s.vocabSvc.GetVocabularyByID(ctx, vocabWordData.VocabID)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.AddWord - get dictionary: %w", err)
	}

	vocabWordData.Native.LangCode = vocab.NativeLang
	vocabWordData.Native.Creator = vocab.UserID

	var nativeWordID uuid.UUID
	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		nativeWords, err := s.dictSvc.GetOrAddWords(ctx, []entityDict.DictWord{vocabWordData.Native})
		if err != nil {
			return fmt.Errorf("add native word in dictionary: %w", err)
		}
		nativeWordID = nativeWords[0].ID

		for i := 0; i < len(vocabWordData.Translates); i++ {
			vocabWordData.Translates[i].LangCode = vocab.TranslateLang
			vocabWordData.Translates[i].Creator = vocab.UserID
		}
		translateWords, err := s.dictSvc.GetOrAddWords(ctx, vocabWordData.Translates)
		if err != nil {
			return fmt.Errorf("add translate word in dictionary: %w", err)
		}
		translateWordIDs := make([]uuid.UUID, 0, len(translateWords))
		for _, word := range translateWords {
			translateWordIDs = append(translateWordIDs, word.ID)
		}

		exampleIDs, err := s.exampleSvc.AddExamples(ctx, vocabWordData.Examples, vocab.NativeLang)
		if err != nil {
			return fmt.Errorf("add example: %w", err)
		}

		err = s.repo.AddWord(ctx, VocabWord{
			ID:           vocabWordData.ID,
			VocabID:      vocabWordData.VocabID,
			NativeID:     nativeWordID,
			TranslateIDs: translateWordIDs,
			ExampleIDs:   exampleIDs,
		})
		if err != nil {
			switch {
			case errors.Is(err, ErrDuplicate):
				return fmt.Errorf("add vocabulary: %w", ErrDuplicate)
			default:
				return fmt.Errorf("add vocabulary: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.AddWord: %w", err)
	}

	vocabularyWord := VocabWord{
		ID:        vocabWordData.ID,
		NativeID:  nativeWordID,
		CreatedAt: vocabWordData.CreatedAt,
		UpdatedAt: vocabWordData.UpdatedAt,
	}

	return vocabularyWord, nil
}

func (s *Service) UpdateWord(ctx context.Context, vocabWordData VocabWordData) (VocabWord, error) {
	vocab, err := s.vocabSvc.GetVocabularyByID(ctx, vocabWordData.VocabID)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - get dictionary: %w", err)
	}

	vocabWordData.Native.LangCode = vocab.NativeLang
	vocabWordData.Native.Creator = vocab.UserID

	nativeWords, err := s.dictSvc.GetOrAddWords(ctx, []entityDict.DictWord{vocabWordData.Native})
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add native word in dictionary: %w", err)
	}
	nativeWordID := nativeWords[0].ID

	for i := 0; i < len(vocabWordData.Translates); i++ {
		vocabWordData.Translates[i].LangCode = vocab.TranslateLang
		vocabWordData.Translates[i].Creator = vocab.UserID
	}
	translateWords, err := s.dictSvc.GetOrAddWords(ctx, vocabWordData.Translates)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add translate word in dictionary: %w", err)
	}
	translateWordIDs := make([]uuid.UUID, 0, len(translateWords))
	for _, word := range translateWords {
		translateWordIDs = append(translateWordIDs, word.ID)
	}

	exampleIDs, err := s.exampleSvc.AddExamples(ctx, vocabWordData.Examples, vocab.NativeLang)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add example: %w", err)
	}

	vocabWord := VocabWord{
		ID:           vocabWordData.ID,
		VocabID:      vocabWordData.VocabID,
		NativeID:     nativeWordID,
		TranslateIDs: translateWordIDs,
		ExampleIDs:   exampleIDs,
		UpdatedAt:    vocabWordData.UpdatedAt,
	}

	err = s.repo.UpdateWord(ctx, vocabWord)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - update vocabulary: %w", err)
	}

	return vocabWord, nil
}

func (s *Service) DeleteWord(ctx context.Context, vocabID, wordID uuid.UUID) error {
	vocabWord := VocabWord{
		ID:      wordID,
		VocabID: vocabID,
	}

	err := s.repo.DeleteWord(ctx, vocabWord)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord - delete word: %w", err)
	}
	return nil
}

func (s *Service) GetRandomWords(ctx context.Context, vocabID uuid.UUID, limit int) ([]VocabWordData, error) {
	vocabWordsData, err := s.repo.GetRandomVocabulary(ctx, vocabID, limit)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, nil
}

func (s *Service) GetWord(ctx context.Context, wordID uuid.UUID) (*VocabWordData, error) {
	vocabWordData, err := s.repo.GetWord(ctx, wordID)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}

	return &vocabWordData, nil
}

func (s *Service) GetWords(ctx context.Context, vocabID uuid.UUID) ([]VocabWordData, error) {
	vocabWordsData, err := s.repo.GetVocabularyWords(ctx, vocabID)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, nil
}
