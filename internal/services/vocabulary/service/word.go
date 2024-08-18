package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/handler"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"
	"github.com/av-ugolkov/lingua-evo/runtime"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type (
	repoWord interface {
		GetWord(ctx context.Context, wordID uuid.UUID) (entity.VocabWordData, error)
		AddWord(ctx context.Context, word entity.VocabWord) (uuid.UUID, error)
		DeleteWord(ctx context.Context, word entity.VocabWord) error
		GetRandomVocabulary(ctx context.Context, vid uuid.UUID, limit int) ([]entity.VocabWordData, error)
		GetVocabularyWords(ctx context.Context, vid uuid.UUID) ([]entity.VocabWordData, error)
		UpdateWord(ctx context.Context, word entity.VocabWord) error
		GetCountWords(ctx context.Context, uid uuid.UUID) (int, error)
	}

	userSvc interface {
		UserCountWord(ctx context.Context, uid uuid.UUID) (int, error)
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

func (s *Service) AddWord(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	userCountWord, err := s.userSvc.UserCountWord(ctx, uid)
	if err != nil {
		return entity.VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord - get count words: %w", err),
			http.StatusInternalServerError, handler.ErrInternal)
	}
	count, err := s.repoVocab.GetCountWords(ctx, uid)
	if err != nil {
		return entity.VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord - get count words: %v", err),
			http.StatusInternalServerError, handler.ErrInternal)
	}

	if count >= userCountWord {
		return entity.VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord: %v", entity.ErrUserWordLimit),
			http.StatusInternalServerError, "You reached word limit")
	}

	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, handler.NewError(fmt.Errorf("word.Service.AddWord - get dictionary: %v", err),
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

		vocabWordData.ID, err = s.repoVocab.AddWord(ctx, entity.VocabWord{
			VocabID:       vocabWordData.VocabID,
			NativeID:      nativeWordID,
			Pronunciation: vocabWordData.Native.Pronunciation,
			TranslateIDs:  translateWordIDs,
			ExampleIDs:    exampleIDs,
		})
		if err != nil {
			switch {
			case errors.Is(err, entity.ErrDuplicate):
				return fmt.Errorf("add vocabulary: %w", entity.ErrDuplicate)
			default:
				return fmt.Errorf("add vocabulary: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.Service.AddWord: %w", err)
	}

	vocabularyWord := entity.VocabWord{
		ID:        vocabWordData.ID,
		NativeID:  nativeWordID,
		CreatedAt: vocabWordData.CreatedAt,
		UpdatedAt: vocabWordData.UpdatedAt,
	}

	return vocabularyWord, nil
}

func (s *Service) UpdateWord(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.Service.UpdateWord - get dictionary: %w", err)
	}

	vocabWordData.Native.LangCode = vocab.NativeLang
	vocabWordData.Native.Creator = vocab.UserID

	nativeWords, err := s.dictSvc.GetOrAddWords(ctx, []entityDict.DictWord{vocabWordData.Native})
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add native word in dictionary: %w", err)
	}
	nativeWordID := nativeWords[0].ID

	for i := 0; i < len(vocabWordData.Translates); i++ {
		vocabWordData.Translates[i].LangCode = vocab.TranslateLang
		vocabWordData.Translates[i].Creator = vocab.UserID
	}
	translateWords, err := s.dictSvc.GetOrAddWords(ctx, vocabWordData.Translates)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add translate word in dictionary: %w", err)
	}
	translateWordIDs := make([]uuid.UUID, 0, len(translateWords))
	for _, word := range translateWords {
		translateWordIDs = append(translateWordIDs, word.ID)
	}

	exampleIDs, err := s.exampleSvc.AddExamples(ctx, vocabWordData.Examples, vocab.NativeLang)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.Service.UpdateWord - add example: %w", err)
	}

	vocabWord := entity.VocabWord{
		ID:            vocabWordData.ID,
		VocabID:       vocabWordData.VocabID,
		NativeID:      nativeWordID,
		Pronunciation: vocabWordData.Native.Pronunciation,
		TranslateIDs:  translateWordIDs,
		ExampleIDs:    exampleIDs,
		UpdatedAt:     vocabWordData.UpdatedAt,
	}

	err = s.repoVocab.UpdateWord(ctx, vocabWord)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.Service.UpdateWord - update vocabulary: %w", err)
	}

	return vocabWord, nil
}

func (s *Service) DeleteWord(ctx context.Context, vid, wid uuid.UUID) error {
	vocabWord := entity.VocabWord{
		ID:      wid,
		VocabID: vid,
	}

	err := s.repoVocab.DeleteWord(ctx, vocabWord)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord - delete word: %w", err)
	}
	return nil
}

func (s *Service) GetRandomWords(ctx context.Context, vid uuid.UUID, limit int) ([]entity.VocabWordData, error) {
	vocabWordsData, err := s.repoVocab.GetRandomVocabulary(ctx, vid, limit)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, nil
}

func (s *Service) GetWord(ctx context.Context, wid uuid.UUID) (*entity.VocabWordData, error) {
	vocabWordData, err := s.repoVocab.GetWord(ctx, wid)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}

	return &vocabWordData, nil
}

func (s *Service) GetWords(ctx context.Context, uid, vid uuid.UUID) ([]entity.VocabWordData, bool, error) {
	_, err := s.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return nil, false, fmt.Errorf("word.Service.GetWords - check access: %w", err)
	}

	vocab, err := s.GetVocabulary(ctx, uid, vid)
	if err != nil {
		return nil, false, fmt.Errorf("word.Service.GetWords - get vocabulary: %w", err)
	}

	editable := vocab.UserID == uid
	if vocab.UserID != uid {
		editable, err = s.VocabularyEditable(ctx, uid, vid)
		if err != nil {
			switch {
			case !errors.Is(err, pgx.ErrNoRows):
				return nil, false, fmt.Errorf("word.Service.GetWords - check access: %w", err)
			}
		}
	}

	vocabWordsData, err := s.repoVocab.GetVocabularyWords(ctx, vid)
	if err != nil {
		return nil, false, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, editable, nil
}

func (s *Service) GetPronunciation(ctx context.Context, uid, vid uuid.UUID, text string) (string, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vid)
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("word.Service.GetPronunciation - get vocabulary: %w", err)
	}
	words, err := s.dictSvc.GetWordsByText(ctx, []entityDict.DictWord{{Text: text, LangCode: vocab.NativeLang}})
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("word.Service.GetPronunciation - get word: %w", err)
	}
	if len(words) == 0 {
		return runtime.EmptyString,
			handler.NewError(fmt.Errorf("word.Service.GetPronunciation - word not found"),
				http.StatusInternalServerError, entity.ErrWordPronunciation.Error())
	}
	word := words[0]
	if word.Pronunciation == runtime.EmptyString {
		return runtime.EmptyString,
			handler.NewError(fmt.Errorf("word.Service.GetPronunciation: %w", entity.ErrWordPronunciation),
				http.StatusInternalServerError, entity.ErrWordPronunciation.Error())
	}
	return word.Pronunciation, nil
}

func (s *Service) CopyWords(ctx context.Context, vid, copyVid uuid.UUID) error {
	vocabWordsData, err := s.repoVocab.GetVocabularyWords(ctx, vid)
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

		_, err = s.repoVocab.AddWord(ctx, entity.VocabWord{
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
