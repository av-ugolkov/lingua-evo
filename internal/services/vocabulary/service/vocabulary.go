package service

import (
	"context"
	"fmt"
	"log/slog"

	"lingua-evo/internal/services/vocabulary/dto"
	"lingua-evo/internal/services/vocabulary/entity"
	entityWord "lingua-evo/internal/services/word/entity"

	"github.com/google/uuid"
)

type (
	repoDict interface {
		AddWord(ctx context.Context, vocabulary entity.Vocabulary) (uuid.UUID, error)
	}

	wordSvc interface {
		AddWord(ctx context.Context, word *entityWord.Word) (uuid.UUID, error)
	}
)

type VocabularySvc struct {
	repo    repoDict
	wordSvc wordSvc
}

func NewService(repo repoDict, wordSvc wordSvc) *VocabularySvc {
	return &VocabularySvc{
		repo:    repo,
		wordSvc: wordSvc,
	}
}

func (s *VocabularySvc) AddWordInVocabulary(
	ctx context.Context,
	vocab *dto.AddWordRq,
) (uuid.UUID, error) {

	word := entityWord.Word{
		Text:          vocab.NativeWord.Text,
		Pronunciation: vocab.NativeWord.Pronunciation,
		LanguageCode:  vocab.NativeWord.LangCode,
	}
	nativeWordID, err := s.wordSvc.AddWord(ctx, &word)
	if err != nil {
		return uuid.Nil, fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add native word in dictionary: %w", err)
	}

	translateWordIDs := make([]uuid.UUID, 0, len(vocab.TanslateWords))
	for _, translateWord := range vocab.TanslateWords {
		word = entityWord.Word{
			Text:          translateWord.Text,
			Pronunciation: translateWord.Pronunciation,
			LanguageCode:  translateWord.LangCode,
		}
		translateWordID, err := s.wordSvc.AddWord(ctx, &word)
		if err != nil {
			return uuid.Nil, fmt.Errorf("vocabulary.service.VocabularuSvc.AddWordInVocabulary - add translate word in dictionary: %w", err)
		}
		translateWordIDs = append(translateWordIDs, translateWordID)
	}
	//TODO сначала получать id слов и остальные данные и потом создается структура
	v := entity.Vocabulary{
		DictionaryId:  vocab.DictionaryID,
		NativeWord:    nativeWordID,
		TranslateWord: translateWordIDs,
		Examples:      []uuid.UUID{},
		Tags:          []uuid.UUID{},
	}

	slog.Info(fmt.Sprintf("vocab: %v", vocab))

	fmt.Print(v)

	vocabID, err := s.repo.AddWord(ctx, v)
	if err != nil {
		return uuid.Nil, err
	}

	return vocabID, nil
}
