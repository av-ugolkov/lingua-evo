package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type repoVocabAdmin interface {
	ChangeVocabTranslationLang(ctx context.Context, vid uuid.UUID, lang string) error
	MoveTranslatedWordsToNewDictionary(ctx context.Context, vid uuid.UUID, lang string) error
	DeleteTranslatedWordsFromOldDictionary(ctx context.Context, vid uuid.UUID) error
}

func (s *Service) ChangeVocabTranslationLang(ctx context.Context, uid uuid.UUID, vid uuid.UUID, lang string) error {
	if ok, _ := s.userSvc.CheckAdmin(ctx, uid); !ok {
		return fmt.Errorf("vocabulary.Service.ChangeVocabTranslationLang: user is not admin")
	}

	err := s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		err := s.repoVocab.ChangeVocabTranslationLang(ctx, vid, lang)
		if err != nil {
			return err
		}

		err = s.repoVocab.MoveTranslatedWordsToNewDictionary(ctx, vid, lang)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("vocabulary.Service.ChangeVocabTranslationLang: %w", err)
	}

	return nil
}
