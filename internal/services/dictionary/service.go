package dictionary

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	errCountDictionary = errors.New("too much dictionaries for user")
)

type (
	repoDict interface {
		AddDictionary(ctx context.Context, dict Dictionary) error
		DeleteDictionary(ctx context.Context, dict Dictionary) error
		GetDictionaryByName(ctx context.Context, dict Dictionary) (uuid.UUID, []uuid.UUID, error)
		GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*Dictionary, error)
	}

	repoVocab interface {
		GetWordsFromDictionary(ctx context.Context, id uuid.UUID, capacity int) ([]string, error)
	}
)

type Service struct {
	repoDict  repoDict
	repoVocab repoVocab
}

func NewService(repoDict repoDict, repoVocab repoVocab) *Service {
	return &Service{
		repoDict:  repoDict,
		repoVocab: repoVocab,
	}
}

func (s *Service) AddDictionary(ctx context.Context, userID, dictID uuid.UUID, name string) (uuid.UUID, error) {
	dictionary := Dictionary{
		ID:     dictID,
		UserID: userID,
		Name:   name,
	}

	dictionaries, err := s.repoDict.GetDictionaries(ctx, dictionary.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.AddDictionary - get count dictionaries: %w", err)
	}

	for _, dict := range dictionaries {
		if dict.Name == dictionary.Name {
			return dict.ID, fmt.Errorf("dictionary.Service.AddDictionary - already have dictionary with same")
		}
	}

	if len(dictionaries) >= 5 {
		return uuid.Nil, fmt.Errorf("dictionary.Service.AddDictionary - %w %v", errCountDictionary, dictionary.UserID)
	}

	err = s.repoDict.AddDictionary(ctx, dictionary)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dictionary.Service.AddDictionary: %w", err)
	}

	return dictionary.ID, nil
}

func (s *Service) DeleteDictionary(ctx context.Context, userID uuid.UUID, name string) error {
	dict := Dictionary{
		UserID: userID,
		Name:   name,
	}

	err := s.repoDict.DeleteDictionary(ctx, dict)
	if err != nil {
		return fmt.Errorf("dictionary.Service.DeleteDictionary: %w", err)
	}
	return nil
}

func (s *Service) GetDictionary(ctx context.Context, userID uuid.UUID, name string) (uuid.UUID, []uuid.UUID, error) {
	dict := Dictionary{
		UserID: userID,
		Name:   name,
	}

	dictID, tags, err := s.repoDict.GetDictionaryByName(ctx, dict)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("dictionary.Service.GetDictionary: %w", err)
	}

	return dictID, tags, nil
}

func (s *Service) GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*Dictionary, error) {
	dictionaries, err := s.repoDict.GetDictionaries(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetDictionaries: %w", err)
	}
	return dictionaries, nil
}
