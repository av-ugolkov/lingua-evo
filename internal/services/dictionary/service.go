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
		Add(ctx context.Context, dict Dictionary) error
		Delete(ctx context.Context, dict Dictionary) error
		GetByName(ctx context.Context, uid uuid.UUID, name string) (Dictionary, error)
		GetByID(ctx context.Context, dictID uuid.UUID) (Dictionary, error)
		GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*Dictionary, error)
		Rename(ctx context.Context, id uuid.UUID, newName string) error
	}

	repoVocab interface {
		GetWordsFromDictionary(ctx context.Context, id uuid.UUID, capacity int) ([]string, error)
	}

	langSvc interface {
		GetLangByCode(ctx context.Context, code string) (string, error)
	}
)

type Service struct {
	repoDict  repoDict
	repoVocab repoVocab
	langSvc   langSvc
}

func NewService(repoDict repoDict, repoVocab repoVocab, langSvc langSvc) *Service {
	return &Service{
		repoDict:  repoDict,
		repoVocab: repoVocab,
		langSvc:   langSvc,
	}
}

func (s *Service) AddDictionary(ctx context.Context, userID, dictID uuid.UUID, name, nativeLangCode, secondLangCode string) (Dictionary, error) {
	dictionary := Dictionary{
		ID:         dictID,
		UserID:     userID,
		Name:       name,
		NativeLang: nativeLangCode,
		SecondLang: secondLangCode,
		Tags:       []string{},
	}

	dictionaries, err := s.repoDict.GetDictionaries(ctx, dictionary.UserID)
	if err != nil {
		return Dictionary{}, fmt.Errorf("dictionary.Service.AddDictionary - get count dictionaries: %w", err)
	}

	for _, dict := range dictionaries {
		if dict.Name == dictionary.Name {
			return Dictionary{}, fmt.Errorf("dictionary.Service.AddDictionary - already have dictionary with same")
		}
	}

	if len(dictionaries) >= 5 {
		return Dictionary{}, fmt.Errorf("dictionary.Service.AddDictionary - %w %v", errCountDictionary, dictionary.UserID)
	}

	err = s.repoDict.Add(ctx, dictionary)
	if err != nil {
		return Dictionary{}, fmt.Errorf("dictionary.Service.AddDictionary: %w", err)
	}

	dictionary.NativeLang, err = s.langSvc.GetLangByCode(ctx, dictionary.NativeLang)
	if err != nil {
		return Dictionary{}, fmt.Errorf("dictionary.Service.AddDictionary - get native lang: %w", err)
	}
	dictionary.SecondLang, err = s.langSvc.GetLangByCode(ctx, dictionary.SecondLang)
	if err != nil {
		return Dictionary{}, fmt.Errorf("dictionary.Service.AddDictionary - get second lang: %w", err)
	}
	return dictionary, nil
}

func (s *Service) DeleteDictionary(ctx context.Context, userID uuid.UUID, name string) error {
	dict := Dictionary{
		UserID: userID,
		Name:   name,
	}

	err := s.repoDict.Delete(ctx, dict)
	if err != nil {
		return fmt.Errorf("dictionary.Service.DeleteDictionary: %w", err)
	}
	return nil
}

func (s *Service) GetDictionary(ctx context.Context, userID uuid.UUID, name string) (Dictionary, error) {
	dict, err := s.repoDict.GetByName(ctx, userID, name)
	if err != nil {
		return dict, fmt.Errorf("dictionary.Service.GetDictionary: %w", err)
	}

	return dict, nil
}

func (s *Service) GetDictByID(ctx context.Context, dictID uuid.UUID) (Dictionary, error) {
	dict, err := s.repoDict.GetByID(ctx, dictID)
	if err != nil {
		return Dictionary{}, fmt.Errorf("dictionary.Service.GetDictionary: %w", err)
	}

	return dict, nil
}

func (s *Service) GetDictionaries(ctx context.Context, userID uuid.UUID) ([]*Dictionary, error) {
	dictionaries, err := s.repoDict.GetDictionaries(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("dictionary.Service.GetDictionaries: %w", err)
	}
	return dictionaries, nil
}

func (s *Service) RenameDictionary(ctx context.Context, id uuid.UUID, newName string) error {
	err := s.repoDict.Rename(ctx, id, newName)
	if err != nil {
		return fmt.Errorf("dictionary.Service.RenameDictionary: %w", err)
	}
	return nil
}
