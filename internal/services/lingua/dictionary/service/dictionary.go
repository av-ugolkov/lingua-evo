package service

import (
	"context"
	"errors"
	"fmt"

	entity "lingua-evo/internal/services/lingua/dictionary"

	"github.com/google/uuid"
)

var (
	errCountDictionary = errors.New("too much dictionaries for user")
)

type (
	repoDict interface {
		AddDictionary(ctx context.Context, dict entity.Dictionary) error
		DeleteDictionary(ctx context.Context, dict entity.Dictionary) error
		GetDictionary(ctx context.Context, dict entity.Dictionary) (uuid.UUID, error)
		GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*entity.Dictionary, error)
	}

	repoVocab interface {
		GetWordsFromDictionary(ctx context.Context, id uuid.UUID, capacity int) ([]string, error)
	}
)

type DictionarySvc struct {
	repoDict  repoDict
	repoVocab repoVocab
}

func NewService(repoDict repoDict, repoVocab repoVocab) *DictionarySvc {
	return &DictionarySvc{
		repoDict:  repoDict,
		repoVocab: repoVocab,
	}
}

func (s *DictionarySvc) AddDictionary(ctx context.Context, userID uuid.UUID, name string) (uuid.UUID, error) {
	dictionary := entity.Dictionary{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
	}

	dictionaries, err := s.repoDict.GetDictionaries(ctx, dictionary.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.service.DictionarySvc.AddDictionary - get count dictionaries: %w", err)
	}

	if len(dictionaries) > 3 {
		return uuid.Nil, fmt.Errorf("dictionary.service.DictionarySvc.AddDictionary - %w %v", errCountDictionary, dictionary.UserID)
	}

	err = s.repoDict.AddDictionary(ctx, dictionary)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.service.DictionarySvc.AddDictionary: %w", err)
	}

	return dictionary.ID, nil
}

func (s *DictionarySvc) DeleteDictionary(ctx context.Context, userID uuid.UUID, name string) error {
	dict := entity.Dictionary{
		UserID: userID,
		Name:   name,
	}

	err := s.repoDict.DeleteDictionary(ctx, dict)
	if err != nil {
		return fmt.Errorf("dictionary.service.DictionarySvc.DeleteDictionary: %w", err)
	}
	return nil
}

func (s *DictionarySvc) GetDictionary(ctx context.Context, userID uuid.UUID, name string, capacity int) (uuid.UUID, []string, error) {
	dict := entity.Dictionary{
		UserID: userID,
		Name:   name,
	}

	dictID, err := s.repoDict.GetDictionary(ctx, dict)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("dictionary.service.DictionarySvc.GetDictionary: %w", err)
	}

	var words []string
	if capacity > 0 {
		words, err = s.repoVocab.GetWordsFromDictionary(ctx, dictID, capacity)
		if err != nil {
			return uuid.Nil, nil, fmt.Errorf("dictionary.service.DictionarySvc.GetDictionary - get words: %w", err)
		}
	}

	return dictID, words, nil
}

func (s *DictionarySvc) GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*entity.Dictionary, error) {
	dictionaries, err := s.repoDict.GetDictionaries(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("dictionary.service.DictionarySvc.GetDictionaries: %w", err)
	}
	return dictionaries, nil
}
