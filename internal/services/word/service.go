package word

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entityVocab "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"
	"github.com/av-ugolkov/lingua-evo/runtime/access"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

type (
	repoWord interface {
		GetWord(ctx context.Context, wordID uuid.UUID) (VocabWordData, error)
		AddWord(ctx context.Context, word VocabWord) error
		DeleteWord(ctx context.Context, word VocabWord) error
		GetRandomVocabulary(ctx context.Context, vid uuid.UUID, limit int) ([]VocabWordData, error)
		GetVocabularyWords(ctx context.Context, vid uuid.UUID) ([]VocabWordData, error)
		UpdateWord(ctx context.Context, word VocabWord) error
		GetCountWords(ctx context.Context, uid uuid.UUID) (int, error)
	}

	userSvc interface {
		UserCountWord(ctx context.Context, uid uuid.UUID) (int, error)
	}

	vocabSvc interface {
		GetVocabulary(ctx context.Context, uid, vid uuid.UUID) (entityVocab.Vocabulary, error)
		GetAccessForUser(ctx context.Context, uid, vid uuid.UUID) (access.Status, error)
		VocabularyEditable(ctx context.Context, uid, vid uuid.UUID) (bool, error)
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
	userSvc    userSvc
	vocabSvc   vocabSvc
	dictSvc    dictSvc
	exampleSvc exampleSvc
}

func NewService(
	tr *transactor.Transactor,
	repo repoWord,
	userSvc userSvc,
	vocabSvc vocabSvc,
	dictSvc dictSvc,
	exampleSvc exampleSvc,
) *Service {
	return &Service{
		tr:         tr,
		repo:       repo,
		userSvc:    userSvc,
		vocabSvc:   vocabSvc,
		dictSvc:    dictSvc,
		exampleSvc: exampleSvc,
	}
}

func (s *Service) AddWord(ctx context.Context, uid uuid.UUID, vocabWordData VocabWordData) (VocabWord, error) {
	userCountWord, err := s.userSvc.UserCountWord(ctx, uid)
	if err != nil {
		return VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord - get count words: %w", err),
			http.StatusInternalServerError, handler.ErrInternal)
	}
	count, err := s.repo.GetCountWords(ctx, uid)
	if err != nil {
		return VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord - get count words: %v", err),
			http.StatusInternalServerError, handler.ErrInternal)
	}

	if count >= userCountWord {
		return VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord: %v", ErrUserWordLimit),
			http.StatusInternalServerError, "You reached word limit")
	}

	vocab, err := s.vocabSvc.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord - get dictionary: %v", err),
			http.StatusInternalServerError, handler.ErrInternal)
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
			VocabID:       vocabWordData.VocabID,
			NativeID:      nativeWordID,
			Pronunciation: vocabWordData.Native.Pronunciation,
			TranslateIDs:  translateWordIDs,
			ExampleIDs:    exampleIDs,
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

func (s *Service) UpdateWord(ctx context.Context, uid uuid.UUID, vocabWordData VocabWordData) (VocabWord, error) {
	vocab, err := s.vocabSvc.GetVocabulary(ctx, uid, vocabWordData.VocabID)
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
		ID:            vocabWordData.ID,
		VocabID:       vocabWordData.VocabID,
		NativeID:      nativeWordID,
		Pronunciation: vocabWordData.Native.Pronunciation,
		TranslateIDs:  translateWordIDs,
		ExampleIDs:    exampleIDs,
		UpdatedAt:     vocabWordData.UpdatedAt,
	}

	err = s.repo.UpdateWord(ctx, vocabWord)
	if err != nil {
		return VocabWord{}, fmt.Errorf("word.Service.UpdateWord - update vocabulary: %w", err)
	}

	return vocabWord, nil
}

func (s *Service) DeleteWord(ctx context.Context, vid, wid uuid.UUID) error {
	vocabWord := VocabWord{
		ID:      wid,
		VocabID: vid,
	}

	err := s.repo.DeleteWord(ctx, vocabWord)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord - delete word: %w", err)
	}
	return nil
}

func (s *Service) GetRandomWords(ctx context.Context, vid uuid.UUID, limit int) ([]VocabWordData, error) {
	vocabWordsData, err := s.repo.GetRandomVocabulary(ctx, vid, limit)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, nil
}

func (s *Service) GetWord(ctx context.Context, wid uuid.UUID) (*VocabWordData, error) {
	vocabWordData, err := s.repo.GetWord(ctx, wid)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}

	return &vocabWordData, nil
}

func (s *Service) GetWords(ctx context.Context, uid, vid uuid.UUID) ([]VocabWordData, bool, error) {
	_, err := s.vocabSvc.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return nil, false, fmt.Errorf("word.Service.GetWords - check access: %w", err)
	}

	vocab, err := s.vocabSvc.GetVocabulary(ctx, uid, vid)
	if err != nil {
		return nil, false, fmt.Errorf("word.Service.GetWords - get vocabulary: %w", err)
	}

	editable := vocab.UserID == uid
	if vocab.UserID != uid {
		editable, err = s.vocabSvc.VocabularyEditable(ctx, uid, vid)
		if err != nil {
			switch {
			case !errors.Is(err, pgx.ErrNoRows):
				return nil, false, fmt.Errorf("word.Service.GetWords - check access: %w", err)
			}
		}
	}

	vocabWordsData, err := s.repo.GetVocabularyWords(ctx, vid)
	if err != nil {
		return nil, false, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, editable, nil
}

func (s *Service) GetPronunciation(ctx context.Context, uid, vid uuid.UUID, text string) (string, error) {
	vocab, err := s.vocabSvc.GetVocabulary(ctx, uid, vid)
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("word.Service.GetPronunciation - get vocabulary: %w", err)
	}
	words, err := s.dictSvc.GetWordsByText(ctx, []entityDict.DictWord{{Text: text, LangCode: vocab.NativeLang}})
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("word.Service.GetPronunciation - get word: %w", err)
	}
	word := words[0]
	if word.Pronunciation == runtime.EmptyString {
		return runtime.EmptyString, fmt.Errorf("word.Service.GetPronunciation: %w", ErrWordPronunciation)
	}
	return word.Pronunciation, nil
}

func (s *Service) CopyWords(ctx context.Context, vid, copyVid uuid.UUID) error {
	vocabWordsData, err := s.repo.GetVocabularyWords(ctx, vid)
	if err != nil {
		return fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	for _, word := range vocabWordsData {
		trIDs := make([]uuid.UUID, 0, len(word.Translates))
		for _, tr := range word.Translates {
			trIDs = append(trIDs, tr.ID)
		}

		exIDs := make([]uuid.UUID, 0, len(word.Examples))
		for _, ex := range word.Examples {
			exIDs = append(exIDs, ex.ID)
		}

		err = s.repo.AddWord(ctx, VocabWord{
			VocabID:       copyVid,
			NativeID:      word.Native.ID,
			Pronunciation: word.Native.Pronunciation,
			TranslateIDs:  trIDs,
			ExampleIDs:    exIDs,
		})
		if err != nil {
			return fmt.Errorf("word.Service.AddWord - add word: %w", err)
		}
	}

	return nil
}
