package word

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entityVocab "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
)

const countWorker = 6

type (
	repoWord interface {
		GetWord(ctx context.Context, vocabID, wordID uuid.UUID) (VocabWord, error)
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
		AddWords(ctx context.Context, words []entityDict.DictWord) ([]uuid.UUID, error)
		UpdateWord(ctx context.Context, word entityDict.DictWord) (uuid.UUID, error)
		GetWords(ctx context.Context, wordIDs []uuid.UUID) ([]entityDict.DictWord, error)
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

func (s *Service) AddWord(ctx context.Context, vocabWord VocabWordData) (VocabWord, error) {
	vocab, err := s.vocabSvc.GetVocabularyByID(ctx, vocabWord.VocabID)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.AddWord - get dictionary: %w", err)
	}

	var nativeWordID uuid.UUID
	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		nativeWordIDs, err := s.dictSvc.AddWords(ctx, []entityDict.DictWord{
			{
				ID:            vocabWord.Native.ID,
				Text:          vocabWord.Native.Text,
				Pronunciation: vocabWord.Native.Pronunciation,
				LangCode:      vocab.NativeLang,
				Creator:       vocab.UserID,
			},
		})
		if err != nil {
			return fmt.Errorf("add native word in dictionary: %w", err)
		}
		nativeWordID = nativeWordIDs[0]

		translateWords := make([]entityDict.DictWord, 0, len(vocabWord.Translates))
		for _, translate := range vocabWord.Translates {
			translateWords = append(translateWords, entityDict.DictWord{
				ID:       translate.ID,
				Text:     translate.Text,
				LangCode: vocab.TranslateLang,
				Creator:  vocab.UserID,
			})
		}
		translateWordIDs, err := s.dictSvc.AddWords(ctx, translateWords)
		if err != nil {
			return fmt.Errorf("add translate word in dictionary: %w", err)
		}

		exampleIDs, err := s.exampleSvc.AddExamples(ctx, vocabWord.Examples, vocab.NativeLang)
		if err != nil {
			return fmt.Errorf("add example: %w", err)
		}

		err = s.repo.AddWord(ctx, VocabWord{
			ID:           vocabWord.ID,
			VocabID:      vocabWord.VocabID,
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
		ID: vocabWord.ID,
	}

	return vocabularyWord, nil
}

func (s *Service) UpdateWord(ctx context.Context, vocabWordData VocabWordData) (VocabWord, error) {
	vocab, err := s.vocabSvc.GetVocabularyByID(ctx, vocabWordData.VocabID)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - get dictionary: %w", err)
	}

	nativeWordID, err := s.dictSvc.UpdateWord(ctx, entityDict.DictWord{
		ID:            vocabWordData.Native.ID,
		Text:          vocabWordData.Native.Text,
		Pronunciation: vocabWordData.Native.Pronunciation,
		LangCode:      vocab.NativeLang,
		Creator:       vocab.UserID,
	})
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(vocabWordData.Translates))
	for _, translate := range vocabWordData.Translates {
		translateWordID, err := s.dictSvc.UpdateWord(ctx, entityDict.DictWord{
			ID:            translate.ID,
			Text:          translate.Text,
			Pronunciation: translate.Pronunciation,
			LangCode:      vocab.TranslateLang,
			Creator:       vocab.UserID,
		})
		if err != nil {
			return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}

	exampleIDs, err := s.exampleSvc.AddExamples(ctx, vocabWordData.Examples, vocab.NativeLang)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add example: %w", err)
	}

	vocabWord := VocabWord{
		ID:           vocabWordData.ID,
		VocabID:      vocabWordData.VocabID,
		NativeID:     vocabWordData.Native.ID,
		TranslateIDs: translateWordIDs,
		ExampleIDs:   exampleIDs,
	}

	if vocabWordData.Native.ID != nativeWordID {
		err = s.repo.DeleteWord(ctx, vocabWord)
		if err != nil {
			return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - delete old word: %w", err)
		}

		err = s.repo.AddWord(ctx, vocabWord)
		if err != nil {
			return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add new word: %w", err)
		}

		return VocabWord{
			ID: vocabWord.ID,
		}, nil
	}
	err = s.repo.UpdateWord(ctx, vocabWord)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - update vocabulary: %w", err)
	}

	return vocabWord, nil
}

func (s *Service) DeleteWord(ctx context.Context, vocabID, nativeWordID uuid.UUID) error {
	vocabWord := VocabWord{
		VocabID:  vocabID,
		NativeID: nativeWordID,
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
			words, err := s.dictSvc.GetWords(ctx, []uuid.UUID{vocabWord.NativeID})
			if err != nil {
				return fmt.Errorf("get words: %w", err)
			}
			if len(words) == 0 {
				return fmt.Errorf("not found word by id [%v]", vocabWord.NativeID)
			}

			translates, err := s.dictSvc.GetWords(ctx, vocabWord.TranslateIDs)
			if err != nil {
				return fmt.Errorf("get translate words: %w", err)
			}

			examples, err := s.exampleSvc.GetExamples(ctx, vocabWord.ExampleIDs)
			if err != nil {
				return fmt.Errorf("get examples: %w", err)
			}

			vocabularyWords = append(vocabularyWords, VocabWordData{
				Native:     words[0],
				Translates: translates,
				Examples:   examples,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabularyWords, nil
}

func (s *Service) GetWord(ctx context.Context, vocabID, wordID uuid.UUID) (*VocabWordData, error) {
	vocabWord, err := s.repo.GetWord(ctx, vocabID, wordID)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}
	var vocabWordData VocabWordData
	var eg errgroup.Group
	eg.Go(func() error {
		words, err := s.dictSvc.GetWords(ctx, []uuid.UUID{wordID})
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
		translateWords, err := s.dictSvc.GetWords(ctx, vocabWord.TranslateIDs)
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
		words, err := s.dictSvc.GetWords(ctx, []uuid.UUID{vocabWord.NativeID})
		if err != nil {
			result <- ResultJob{err: fmt.Errorf("get words: %w", err)}
			return
		}
		if len(words) == 0 {
			result <- ResultJob{err: fmt.Errorf("not found word by id [%v]", vocabWord.NativeID)}
			return
		}

		translates, err := s.dictSvc.GetWords(ctx, vocabWord.TranslateIDs)
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
			},
		}
	}

	stopCh <- struct{}{}
}
