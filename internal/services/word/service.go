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
	"github.com/av-ugolkov/lingua-evo/internal/services/word/model"
)

const countWorker = 6

type (
	repoWord interface {
		GetWord(ctx context.Context, dictID, wordID uuid.UUID) (Word, error)
		AddWord(ctx context.Context, word Word) error
		DeleteWord(ctx context.Context, word Word) error
		GetRandomVocabulary(ctx context.Context, dictID uuid.UUID, limit int) ([]Word, error)
		GetVocabulary(ctx context.Context, dictID uuid.UUID) ([]Word, error)
		UpdateWord(ctx context.Context, word Word) error
	}

	vocabSvc interface {
		GetVocabularyByID(ctx context.Context, vocabID uuid.UUID) (entityVocab.Vocabulary, error)
	}

	exampleSvc interface {
		AddExamples(ctx context.Context, examples []entityExample.Example, langCode string) ([]uuid.UUID, error)
		GetExamples(ctx context.Context, exampleIDs []uuid.UUID) ([]entityExample.Example, error)
	}

	dictSvc interface {
		AddWords(ctx context.Context, words []entityDict.Word) ([]uuid.UUID, error)
		UpdateWord(ctx context.Context, word entityDict.Word) (uuid.UUID, error)
		GetWords(ctx context.Context, wordIDs []uuid.UUID) ([]entityDict.Word, error)
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

func (s *Service) AddWord(ctx context.Context, vocabWord VocabWord) (VocabularyWord, error) {
	vocab, err := s.vocabSvc.GetVocabularyByID(ctx, vocabWord.VocabID)
	if err != nil {
		return VocabularyWord{}, fmt.Errorf("word.Service.AddWord - get dictionary: %w", err)
	}

	var nativeWordID uuid.UUID
	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		nativeWordIDs, err := s.dictSvc.AddWords(ctx, []entityDict.Word{
			{
				ID:            vocabWord.NativeWord.ID,
				Text:          vocabWord.NativeWord.Text,
				Pronunciation: vocabWord.NativeWord.Pronunciation,
				LangCode:      vocab.NativeLang,
			},
		})
		if err != nil {
			return fmt.Errorf("add native word in dictionary: %w", err)
		}
		nativeWordID = nativeWordIDs[0]

		translateWords := make([]entityDict.Word, 0, len(vocabWord.TranslateWords))
		for _, translateWord := range vocabWord.TranslateWords {
			translateWords = append(translateWords, entityDict.Word{
				ID:       translateWord.ID,
				Text:     translateWord.Text,
				LangCode: vocab.TranslateLang,
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

		err = s.repo.AddWord(ctx, Word{
			ID:             vocabWord.WordID,
			VocabID:        vocabWord.VocabID,
			NativeID:       nativeWordID,
			TranslateWords: translateWordIDs,
			Examples:       exampleIDs,
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
		return VocabularyWord{}, fmt.Errorf("word.Service.AddWord: %w", err)
	}

	vocabularyWord := VocabularyWord{
		ID: nativeWordID,
	}

	return vocabularyWord, nil
}

func (s *Service) UpdateWord(ctx context.Context,
	vocabID uuid.UUID,
	wordID uuid.UUID,
	nativeWord model.VocabWord,
	translateWords Words,
	examples []entityExample.Example) (VocabularyWord, error) {
	vocab, err := s.vocabSvc.GetVocabularyByID(ctx, vocabID)
	if err != nil {
		return VocabularyWord{}, fmt.Errorf("word.Service.UpdateWord - get dictionary: %w", err)
	}

	nativeWordID, err := s.dictSvc.UpdateWord(ctx, entityDict.Word{
		ID:            uuid.New(),
		Text:          nativeWord.Text,
		LangCode:      vocab.NativeLang,
		Pronunciation: nativeWord.Pronunciation,
	})
	if err != nil {
		return VocabularyWord{}, fmt.Errorf("word.Service.UpdateWord - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(translateWords))
	for _, translateWord := range translateWords {
		translateWordID, err := s.dictSvc.UpdateWord(ctx, entityDict.Word{
			ID:            uuid.New(),
			Text:          translateWord.Text,
			LangCode:      vocab.TranslateLang,
			Pronunciation: translateWord.Pronunciation,
		})
		if err != nil {
			return VocabularyWord{}, fmt.Errorf("word.Service.UpdateWord - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}

	exampleIDs, err := s.exampleSvc.AddExamples(ctx, examples, vocab.NativeLang)
	if err != nil {
		return VocabularyWord{}, fmt.Errorf("word.Service.UpdateWord - add example: %w", err)
	}

	vocabulary := Word{
		ID:             vocabID,
		NativeID:       nativeWordID,
		TranslateWords: translateWordIDs,
		Examples:       exampleIDs,
	}

	if wordID != nativeWordID {
		err = s.repo.DeleteWord(ctx, Word{ID: vocabID, NativeID: wordID})
		if err != nil {
			return VocabularyWord{}, fmt.Errorf("word.Service.UpdateWord - delete old word: %w", err)
		}
		err = s.repo.AddWord(ctx, vocabulary)
		if err != nil {
			return VocabularyWord{}, fmt.Errorf("word.Service.UpdateWord - add new word: %w", err)
		}
		vocabExamples := make([]string, 0, len(examples))
		for _, example := range examples {
			vocabExamples = append(vocabExamples, example.Text)
		}
		return VocabularyWord{
			ID:             nativeWordID,
			NativeWord:     nativeWord,
			TranslateWords: translateWords.GetValues(),
			Examples:       vocabExamples,
		}, nil
	}
	err = s.repo.UpdateWord(ctx, vocabulary)
	if err != nil {
		return VocabularyWord{}, fmt.Errorf("word.Service.UpdateWord - update vocabulary: %w", err)
	}
	vocabExamples := make([]string, 0, len(examples))
	for _, example := range examples {
		vocabExamples = append(vocabExamples, example.Text)
	}
	return VocabularyWord{
		ID:             nativeWordID,
		NativeWord:     nativeWord,
		TranslateWords: translateWords.GetValues(),
		Examples:       vocabExamples,
	}, nil
}

func (s *Service) DeleteWord(ctx context.Context, vocabID, nativeWordID uuid.UUID) error {
	vocabWord := Word{
		ID:       vocabID,
		NativeID: nativeWordID,
	}

	err := s.repo.DeleteWord(ctx, vocabWord)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord - delete word: %w", err)
	}
	return nil
}

func (s *Service) GetRandomWords(ctx context.Context, vocabID uuid.UUID, limit int) ([]VocabularyWord, error) {
	vocabularies, err := s.repo.GetRandomVocabulary(ctx, vocabID, limit)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	vocabularyWords := make([]VocabularyWord, 0, len(vocabularies))

	var eg errgroup.Group
	eg.SetLimit(10)
	for _, vocab := range vocabularies {
		vocab := vocab
		eg.Go(func() error {
			words, err := s.dictSvc.GetWords(ctx, []uuid.UUID{vocab.NativeID})
			if err != nil {
				return fmt.Errorf("get words: %w", err)
			}
			if len(words) == 0 {
				return fmt.Errorf("not found word by id [%v]", vocab.NativeID)
			}

			translateWords, err := s.dictSvc.GetWords(ctx, vocab.TranslateWords)
			if err != nil {
				return fmt.Errorf("get translate words: %w", err)
			}
			translates := make([]string, 0, len(translateWords))
			for _, word := range translateWords {
				translates = append(translates, word.Text)
			}

			examples, err := s.exampleSvc.GetExamples(ctx, vocab.Examples)
			if err != nil {
				return fmt.Errorf("get examples: %w", err)
			}
			examplesStr := make([]string, 0, len(examples))
			for _, example := range examples {
				examplesStr = append(examplesStr, example.Text)
			}

			vocabularyWords = append(vocabularyWords, VocabularyWord{
				NativeWord: model.VocabWord{
					Text:          words[0].Text,
					Pronunciation: words[0].Pronunciation,
				},
				TranslateWords: translates,
				Examples:       examplesStr,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabularyWords, nil
}

func (s *Service) GetWord(ctx context.Context, vocabID, wordID uuid.UUID) (*VocabularyWord, error) {
	word, err := s.repo.GetWord(ctx, vocabID, wordID)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}
	var vocab VocabularyWord
	var eg errgroup.Group
	eg.Go(func() error {
		words, err := s.dictSvc.GetWords(ctx, []uuid.UUID{wordID})
		if err != nil {
			return fmt.Errorf("get words: %w", err)
		}
		if len(words) == 0 {
			return fmt.Errorf("not found word by id [%v]", wordID)
		}

		vocab.NativeWord = model.VocabWord{
			Text:          words[0].Text,
			Pronunciation: words[0].Pronunciation,
		}

		return nil
	})
	eg.Go(func() error {
		translateWords, err := s.dictSvc.GetWords(ctx, word.TranslateWords)
		if err != nil {
			return fmt.Errorf("get translate words: %w", err)
		}
		vocab.TranslateWords = make([]string, 0, len(translateWords))
		for _, word := range translateWords {
			vocab.TranslateWords = append(vocab.TranslateWords, word.Text)
		}
		return nil
	})
	eg.Go(func() error {
		examples, err := s.exampleSvc.GetExamples(ctx, word.Examples)
		if err != nil {
			return fmt.Errorf("get examples: %w", err)
		}
		vocab.Examples = make([]string, 0, len(examples))
		for _, example := range examples {
			vocab.Examples = append(vocab.Examples, example.Text)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get data: %w", err)
	}

	vocab.ID = wordID

	return &vocab, nil
}

type ResultJob struct {
	value VocabularyWord
	err   error
}

func (s *Service) GetWords(ctx context.Context, vocabID uuid.UUID) ([]VocabularyWord, error) {
	vocabWords, err := s.repo.GetVocabulary(ctx, vocabID)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	vocabularyWords := make([]VocabularyWord, 0, len(vocabWords))

	data := make(chan Word, countWorker)
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
	inData <-chan Word,
	result chan<- ResultJob,
	stopCh chan<- struct{}) {
	for vocab := range inData {
		words, err := s.dictSvc.GetWords(ctx, []uuid.UUID{vocab.NativeID})
		if err != nil {
			result <- ResultJob{err: fmt.Errorf("get words: %w", err)}
			return
		}
		if len(words) == 0 {
			result <- ResultJob{err: fmt.Errorf("not found word by id [%v]", vocab.NativeID)}
			return
		}

		translateWords, err := s.dictSvc.GetWords(ctx, vocab.TranslateWords)
		if err != nil {
			result <- ResultJob{err: fmt.Errorf("get translate words: %w", err)}
			return
		}
		translates := make([]string, 0, len(translateWords))
		for _, word := range translateWords {
			translates = append(translates, word.Text)
		}

		examples, err := s.exampleSvc.GetExamples(ctx, vocab.Examples)
		if err != nil {
			result <- ResultJob{err: fmt.Errorf("get examples: %w", err)}
			return
		}
		examplesStr := make([]string, 0, len(examples))
		for _, example := range examples {
			examplesStr = append(examplesStr, example.Text)
		}

		result <- ResultJob{
			value: VocabularyWord{
				ID: words[0].ID,
				NativeWord: model.VocabWord{
					Text:          words[0].Text,
					Pronunciation: words[0].Pronunciation,
				},
				TranslateWords: translates,
				Examples:       examplesStr,
			},
		}
	}

	stopCh <- struct{}{}
}
