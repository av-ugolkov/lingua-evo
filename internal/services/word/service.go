package word

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"sync"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entityVocab "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
)

const countWorker = 6

type (
	repoWord interface {
		GetWord(ctx context.Context, wordID uuid.UUID) (VocabWord, error)
		AddWord(ctx context.Context, word VocabWord) error
		DeleteWord(ctx context.Context, word VocabWord) error
		GetRandomVocabulary(ctx context.Context, vocabID uuid.UUID, limit int) ([]VocabWord, error)
		GetVocabulary(ctx context.Context, vocabID uuid.UUID) ([]VocabWord, error)
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
	vocabWords, err := s.repo.GetRandomVocabulary(ctx, vocabID, limit)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	vocabularyWords := make([]VocabWordData, 0, len(vocabWords))

	var eg errgroup.Group
	eg.SetLimit(10)
	for _, vocabWord := range vocabWords {
		eg.Go(func() error {
			words, err := s.dictSvc.GetWordsByID(ctx, []uuid.UUID{vocabWord.NativeID})
			if err != nil {
				return fmt.Errorf("get words: %w", err)
			}
			if len(words) == 0 {
				return fmt.Errorf("not found word by id [%v]", vocabWord.NativeID)
			}

			translates, err := s.dictSvc.GetWordsByID(ctx, vocabWord.TranslateIDs)
			if err != nil {
				return fmt.Errorf("get translate words: %w", err)
			}

			vocabularyWords = append(vocabularyWords, VocabWordData{
				Native:     words[0],
				Translates: translates,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabularyWords, nil
}

func (s *Service) GetWord(ctx context.Context, wordID uuid.UUID) (*VocabWordData, error) {
	vocabWord, err := s.repo.GetWord(ctx, wordID)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}
	var vocabWordData VocabWordData
	var eg errgroup.Group
	eg.Go(func() error {
		words, err := s.dictSvc.GetWordsByID(ctx, []uuid.UUID{vocabWord.NativeID})
		if err != nil {
			return fmt.Errorf("get words: %w", err)
		}
		if len(words) == 0 {
			return fmt.Errorf("not found word by id [%v]", wordID)
		}

		vocabWordData.Native = words[0]

		return nil
	})
	eg.Go(func() error {
		translateWords, err := s.dictSvc.GetWordsByID(ctx, vocabWord.TranslateIDs)
		if err != nil {
			return fmt.Errorf("get translate words: %w", err)
		}
		vocabWordData.Translates = translateWords
		return nil
	})
	eg.Go(func() error {
		examples, err := s.exampleSvc.GetExamples(ctx, vocabWord.ExampleIDs)
		if err != nil {
			return fmt.Errorf("get examples: %w", err)
		}
		vocabWordData.Examples = examples
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get data: %w", err)
	}

	vocabWordData.ID = wordID

	return &vocabWordData, nil
}

type ResultJob struct {
	value VocabWordData
	err   error
}

func (s *Service) GetWords(ctx context.Context, vocabID uuid.UUID) ([]VocabWordData, error) {
	vocabWords, err := s.repo.GetVocabulary(ctx, vocabID)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	vocabularyWords := make([]VocabWordData, 0, len(vocabWords))

	data := make(chan VocabWord, countWorker)
	result := make(chan ResultJob, countWorker)
	stopChan := make(chan struct{}, countWorker)
	defer close(stopChan)

	var wg sync.WaitGroup
	wg.Add(countWorker)
	go func() {
		wg.Wait()
		close(result)
	}()

	for w := 0; w < countWorker; w++ {
		go s.workerForGetWord(ctx, data, result, stopChan)
	}

	go func() {
		defer close(data)
		for _, vocab := range vocabWords {
			data <- vocab
		}
	}()

loop:
	for {
		select {
		case res, ok := <-result:
			if !ok {
				break loop
			}
			if res.err != nil {
				return nil, fmt.Errorf("word.Service.GetWords: %w", err)
			}
			vocabularyWords = append(vocabularyWords, res.value)
		case <-stopChan:
			wg.Done()
		}
	}

	return vocabularyWords, nil
}

func (s *Service) workerForGetWord(
	ctx context.Context,
	inData <-chan VocabWord,
	result chan<- ResultJob,
	stopCh chan<- struct{}) {
	for vocabWord := range inData {
		words, err := s.dictSvc.GetWordsByID(ctx, []uuid.UUID{vocabWord.NativeID})
		if err != nil {
			result <- ResultJob{err: fmt.Errorf("get words: %w", err)}
			return
		}
		if len(words) == 0 {
			result <- ResultJob{err: fmt.Errorf("not found word by id [%v]", vocabWord.NativeID)}
			return
		}

		translates, err := s.dictSvc.GetWordsByID(ctx, vocabWord.TranslateIDs)
		if err != nil {
			result <- ResultJob{err: fmt.Errorf("get translate words: %w", err)}
			return
		}

		examples, err := s.exampleSvc.GetExamples(ctx, vocabWord.ExampleIDs)
		if err != nil {
			result <- ResultJob{err: fmt.Errorf("get examples: %w", err)}
			return
		}

		result <- ResultJob{
			value: VocabWordData{
				ID:         vocabWord.ID,
				Native:     words[0],
				Translates: translates,
				Examples:   examples,
				CreatedAt:  vocabWord.CreatedAt,
				UpdatedAt:  vocabWord.UpdatedAt,
			},
		}
	}

	stopCh <- struct{}{}
}
