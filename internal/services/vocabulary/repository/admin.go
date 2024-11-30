package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *VocabRepo) ChangeVocabTranslationLang(ctx context.Context, vid uuid.UUID, lang string) error {
	const query = `UPDATE vocabulary SET translate_lang=$2 WHERE id=$1;`

	_, err := r.tr.Exec(ctx, query, vid, lang)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.ChangeVocabTranslationLang: %v", err)
	}

	return nil
}

func (r *VocabRepo) MoveTranslatedWordsToNewDictionary(ctx context.Context, vid uuid.UUID, lang string) error {

	return nil
}

func (r *VocabRepo) DeleteTranslatedWordsFromOldDictionary(ctx context.Context, vid uuid.UUID) error {
	return nil
}
