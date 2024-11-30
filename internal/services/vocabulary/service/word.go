package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/pkg/msg-error"
	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityEvents "github.com/av-ugolkov/lingua-evo/internal/services/events"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
)

type (
	repoWord interface {
		GetWord(ctx context.Context, wordID uuid.UUID, nativeLang, translateLang string) (entity.VocabWordData, error)
		AddWord(ctx context.Context, word entity.VocabWord) (uuid.UUID, error)
		DeleteWord(ctx context.Context, word entity.VocabWord) error
		GetRandomVocabulary(ctx context.Context, vid uuid.UUID, limit int) ([]entity.VocabWordData, error)
		GetVocabWords(ctx context.Context, vid uuid.UUID) ([]entity.VocabWordData, error)
		GetVocabSeveralWords(ctx context.Context, vid uuid.UUID, count int, nativeLang, translateLang string) ([]entity.VocabWordData, error)
		UpdateWord(ctx context.Context, word entity.VocabWord) error
		UpdateWordText(ctx context.Context, word entity.VocabWord) error
		UpdateWordPronunciation(ctx context.Context, word entity.VocabWord) error
		UpdateWordDefinition(ctx context.Context, word entity.VocabWord) error
		UpdateWordTranslates(ctx context.Context, word entity.VocabWord) error
		UpdateWordExamples(ctx context.Context, word entity.VocabWord) error
		GetCountWords(ctx context.Context, uid uuid.UUID) (int, error)
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

	eventsSvc interface {
		AddEvent(ctx context.Context, event entityEvents.Event) (uuid.UUID, error)
		AsyncAddEvent(event entityEvents.Event)
	}
)

func (s *Service) AddWord(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, msgerr.New(fmt.Errorf("word.Service.AddWord: %v", err), msgerr.ErrMsgInternal)
	}

	vocabWordData.Native.LangCode = vocab.NativeLang
	vocabWordData.Native.Creator = vocab.UserID

	var nativeWordID uuid.UUID
	err = s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		nativeWords, err := s.dictSvc.GetOrAddWords(ctx, []entityDict.DictWord{{
			ID:            vocabWordData.Native.ID,
			Text:          vocabWordData.Native.Text,
			Pronunciation: vocabWordData.Native.Pronunciation,
			LangCode:      vocabWordData.Native.LangCode,
			Creator:       vocabWordData.Native.Creator,
		}})
		if err != nil {
			return fmt.Errorf("add native word in dictionary: %w", err)
		}
		nativeWordID = nativeWords[0].ID

		translates := make([]entityDict.DictWord, 0, len(vocabWordData.Translates))
		for _, tr := range vocabWordData.Translates {
			translates = append(translates, entityDict.DictWord{
				Text:     tr.Text,
				LangCode: vocab.TranslateLang,
				Creator:  vocabWordData.Native.Creator,
			})
		}
		translateWords, err := s.dictSvc.GetOrAddWords(ctx, translates)
		if err != nil {
			return fmt.Errorf("add translate word in dictionary: %w", err)
		}
		translateWordIDs := make([]uuid.UUID, 0, len(translateWords))
		for _, word := range translateWords {
			translateWordIDs = append(translateWordIDs, word.ID)
		}

		examples := make([]entityExample.Example, 0, len(vocabWordData.Examples))
		for _, ex := range vocabWordData.Examples {
			examples = append(examples, entityExample.Example{
				Text: ex.Text,
			})
		}
		exampleIDs, err := s.exampleSvc.AddExamples(ctx, examples, vocab.NativeLang)
		if err != nil {
			return fmt.Errorf("add example: %w", err)
		}

		vocabWordData.ID, err = s.repoVocab.AddWord(ctx, entity.VocabWord{
			VocabID:       vocabWordData.VocabID,
			NativeID:      nativeWordID,
			Pronunciation: vocabWordData.Native.Pronunciation,
			Definition:    vocabWordData.Definition,
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

	now := time.Now().UTC()
	vocabularyWord := entity.VocabWord{
		ID:        vocabWordData.ID,
		VocabID:   vocabWordData.VocabID,
		NativeID:  nativeWordID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.eventsSvc.AsyncAddEvent(entityEvents.Event{
		User: entityEvents.UserData{ID: uid},
		Type: entityEvents.VocabWordCreated,
		Payload: entityEvents.PayloadDataVocab{
			DictWordID: &nativeWordID,
			DictWord:   vocabWordData.Native.Text,
			VocabID:    &vocabWordData.VocabID,
			VocabTitle: vocab.Name,
		},
	})

	return vocabularyWord, nil
}

func (s *Service) UpdateWordText(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordText: %w", err)
	}

	vocabWordData.Native.LangCode = vocab.NativeLang
	vocabWordData.Native.Creator = vocab.UserID

	nativeWords, err := s.dictSvc.GetOrAddWords(ctx, []entityDict.DictWord{{
		ID:            vocabWordData.Native.ID,
		Text:          vocabWordData.Native.Text,
		Pronunciation: vocabWordData.Native.Pronunciation,
		LangCode:      vocabWordData.Native.LangCode,
		Creator:       vocabWordData.Native.Creator,
	}})
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordText: %w", err)
	}
	nativeWordID := nativeWords[0].ID

	vocabWord := entity.VocabWord{
		ID:       vocabWordData.ID,
		VocabID:  vocabWordData.VocabID,
		NativeID: nativeWordID,
	}

	err = s.repoVocab.UpdateWordText(ctx, vocabWord)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordText: %w", err)
	}

	eventType := entityEvents.VocabWordUpdated
	if vocabWordData.Native.ID != nativeWordID {
		eventType = entityEvents.VocabWordRenamed
	}
	s.eventsSvc.AsyncAddEvent(entityEvents.Event{
		User: entityEvents.UserData{ID: uid},
		Type: eventType,
		Payload: entityEvents.PayloadDataVocab{
			DictWordID: &nativeWordID,
			DictWord:   vocabWordData.Native.Text,
			VocabID:    &vocabWord.VocabID,
			VocabTitle: vocab.Name,
		},
	})

	return vocabWord, nil
}

func (s *Service) UpdateWordPronunciation(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordPronunciation: %w", err)
	}

	vocabWord := entity.VocabWord{
		ID:            vocabWordData.ID,
		VocabID:       vocabWordData.VocabID,
		Pronunciation: vocabWordData.Native.Pronunciation,
	}

	err = s.repoVocab.UpdateWordPronunciation(ctx, vocabWord)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordPronunciation: %w", err)
	}

	s.eventsSvc.AsyncAddEvent(entityEvents.Event{
		User: entityEvents.UserData{ID: uid},
		Type: entityEvents.VocabWordUpdated,
		Payload: entityEvents.PayloadDataVocab{
			DictWordID: &vocabWordData.Native.ID,
			DictWord:   vocabWordData.Native.Text,
			VocabID:    &vocabWord.VocabID,
			VocabTitle: vocab.Name,
		},
	})

	return vocabWord, nil
}

func (s *Service) UpdateWordDefinition(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordDefinition: %w", err)
	}

	vocabWord := entity.VocabWord{
		ID:         vocabWordData.ID,
		VocabID:    vocabWordData.VocabID,
		Definition: vocabWordData.Definition,
	}

	err = s.repoVocab.UpdateWordDefinition(ctx, vocabWord)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordDefinition: %w", err)
	}

	vocabWordData, err = s.repoVocab.GetWord(ctx, vocabWord.ID, vocab.NativeLang, vocab.TranslateLang)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordDefinition: %w", err)
	}

	s.eventsSvc.AsyncAddEvent(entityEvents.Event{
		User: entityEvents.UserData{ID: uid},
		Type: entityEvents.VocabWordUpdated,
		Payload: entityEvents.PayloadDataVocab{
			DictWordID: &vocabWordData.Native.ID,
			DictWord:   vocabWordData.Native.Text,
			VocabID:    &vocabWord.VocabID,
			VocabTitle: vocab.Name,
		},
	})

	return vocabWord, nil
}

func (s *Service) UpdateWordTranslates(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordTranslates: %w", err)
	}

	translates := make([]entityDict.DictWord, 0, len(vocabWordData.Translates))
	for _, tr := range vocabWordData.Translates {
		translates = append(translates, entityDict.DictWord{
			Text:     tr.Text,
			LangCode: vocab.TranslateLang,
			Creator:  vocab.UserID,
		})
	}
	translateWords, err := s.dictSvc.GetOrAddWords(ctx, translates)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordTranslates: %w", err)
	}
	translateWordIDs := make([]uuid.UUID, 0, len(translateWords))
	for _, word := range translateWords {
		translateWordIDs = append(translateWordIDs, word.ID)
	}

	vocabWord := entity.VocabWord{
		ID:           vocabWordData.ID,
		VocabID:      vocabWordData.VocabID,
		TranslateIDs: translateWordIDs,
	}

	err = s.repoVocab.UpdateWordTranslates(ctx, vocabWord)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordTranslates: %w", err)
	}

	vocabWordData, err = s.repoVocab.GetWord(ctx, vocabWord.ID, vocab.NativeLang, vocab.TranslateLang)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordDefinition: %w", err)
	}

	s.eventsSvc.AsyncAddEvent(entityEvents.Event{
		User: entityEvents.UserData{ID: uid},
		Type: entityEvents.VocabWordUpdated,
		Payload: entityEvents.PayloadDataVocab{
			DictWordID: &vocabWordData.Native.ID,
			DictWord:   vocabWordData.Native.Text,
			VocabID:    &vocabWord.VocabID,
			VocabTitle: vocab.Name,
		},
	})

	return vocabWord, nil
}

func (s *Service) UpdateWordExamples(ctx context.Context, uid uuid.UUID, vocabWordData entity.VocabWordData) (entity.VocabWord, error) {
	vocab, err := s.GetVocabulary(ctx, uid, vocabWordData.VocabID)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordExamples: %w", err)
	}

	examples := make([]entityExample.Example, 0, len(vocabWordData.Examples))
	for _, ex := range vocabWordData.Examples {
		examples = append(examples, entityExample.Example{
			Text: ex.Text,
		})
	}
	exampleIDs, err := s.exampleSvc.AddExamples(ctx, examples, vocab.NativeLang)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordExamples: %w", err)
	}

	vocabWord := entity.VocabWord{
		ID:         vocabWordData.ID,
		VocabID:    vocabWordData.VocabID,
		ExampleIDs: exampleIDs,
	}

	err = s.repoVocab.UpdateWordExamples(ctx, vocabWord)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordExamples: %w", err)
	}

	vocabWordData, err = s.repoVocab.GetWord(ctx, vocabWord.ID, vocab.NativeLang, vocab.TranslateLang)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("service.vocabulary.Service.UpdateWordDefinition: %w", err)
	}

	s.eventsSvc.AsyncAddEvent(entityEvents.Event{
		User: entityEvents.UserData{ID: uid},
		Type: entityEvents.VocabWordUpdated,
		Payload: entityEvents.PayloadDataVocab{
			DictWordID: &vocabWordData.Native.ID,
			DictWord:   vocabWordData.Native.Text,
			VocabID:    &vocabWord.VocabID,
			VocabTitle: vocab.Name,
		},
	})

	return vocabWord, nil
}

func (s *Service) DeleteWord(ctx context.Context, uid, vid, wid uuid.UUID) error {
	vocabWord := entity.VocabWord{
		ID:      wid,
		VocabID: vid,
	}

	vocabWordTemp, err := s.GetWord(ctx, vid, wid)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord: %w", err)
	}

	vocab, err := s.GetVocabulary(ctx, uid, vid)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord: %w", err)
	}

	err = s.repoVocab.DeleteWord(ctx, vocabWord)
	if err != nil {
		return fmt.Errorf("word.Service.DeleteWord - delete word: %w", err)
	}

	s.eventsSvc.AsyncAddEvent(entityEvents.Event{
		User: entityEvents.UserData{ID: uid},
		Type: entityEvents.VocabWordDeleted,
		Payload: entityEvents.PayloadDataVocab{
			DictWordID: &vocabWordTemp.Native.ID,
			DictWord:   vocabWordTemp.Native.Text,
			VocabID:    &vid,
			VocabTitle: vocab.Name,
		},
	})

	return nil
}

func (s *Service) GetRandomWords(ctx context.Context, vid uuid.UUID, limit int) ([]entity.VocabWordData, error) {
	vocabWordsData, err := s.repoVocab.GetRandomVocabulary(ctx, vid, limit)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, nil
}

func (s *Service) GetWord(ctx context.Context, vid, wid uuid.UUID) (*entity.VocabWordData, error) {
	vocab, err := s.repoVocab.GetVocab(ctx, vid)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord - get vocab: %w", err)
	}

	vocabWordData, err := s.repoVocab.GetWord(ctx, wid, vocab.NativeLang, vocab.TranslateLang)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWord: %w", err)
	}

	return &vocabWordData, nil
}

func (s *Service) GetWords(ctx context.Context, uid, vid uuid.UUID) ([]entity.VocabWordData, error) {
	_, err := s.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - check access: %w", err)
	}

	vocabWordsData, err := s.repoVocab.GetVocabWords(ctx, vid)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}
	return vocabWordsData, nil
}

func (s *Service) GetSeveralWords(ctx context.Context, uid, vid uuid.UUID, count int) ([]entity.VocabWordData, error) {
	_, err := s.GetAccessForUser(ctx, uid, vid)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - check access: %w", err)
	}

	vocab, err := s.repoVocab.GetVocab(ctx, vid)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get vocab: %w", err)
	}

	vocabWordsData, err := s.repoVocab.GetVocabSeveralWords(ctx, vid, count, vocab.NativeLang, vocab.TranslateLang)
	if err != nil {
		return nil, fmt.Errorf("word.Service.GetWords - get words: %w", err)
	}

	return vocabWordsData, nil
}

func (s *Service) CopyWords(ctx context.Context, vid, copyVid uuid.UUID) error {
	vocabWordsData, err := s.repoVocab.GetVocabWords(ctx, vid)
	if err != nil {
		return fmt.Errorf("word.Service.GetWords: %w", err)
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
