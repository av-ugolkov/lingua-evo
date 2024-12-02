package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type repoVocabAdmin interface {
	ChangeVocabTranslationLang(ctx context.Context, vid uuid.UUID, lang string) error
	MoveTranslatedWordsToNewDictionary(ctx context.Context, vid uuid.UUID, oldLang, newLang string) error
	UpdateVocabTranslatedIDs(ctx context.Context, vid uuid.UUID, oldLang, newLang string) error
	DeleteTranslatedWordsFromOldDictionary(ctx context.Context, oldLang string) error
}

func (s *Service) ChangeVocabTranslationLang(ctx context.Context, uid uuid.UUID, vid uuid.UUID, newLang string) error {
	if ok, _ := s.userSvc.CheckAdmin(ctx, uid); !ok {
		return fmt.Errorf("vocabulary.Service.ChangeVocabTranslationLang: user is not admin")
	}

	err := s.tr.CreateTransaction(ctx, func(ctx context.Context) error {
		vocab, err := s.repoVocab.GetVocab(ctx, vid)
		if err != nil {
			return err
		}
		oldLang := vocab.TranslateLang

		err = s.repoVocab.ChangeVocabTranslationLang(ctx, vid, newLang)
		if err != nil {
			return err
		}

		err = s.repoVocab.MoveTranslatedWordsToNewDictionary(ctx, vid, oldLang, newLang)
		if err != nil {
			return err
		}

		err = s.repoVocab.UpdateVocabTranslatedIDs(ctx, vid, oldLang, newLang)
		if err != nil {
			return err
		}

		err = s.repoVocab.DeleteTranslatedWordsFromOldDictionary(ctx, oldLang)
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
