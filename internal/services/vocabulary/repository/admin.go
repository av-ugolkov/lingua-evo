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

func (r *VocabRepo) MoveTranslatedWordsToNewDictionary(ctx context.Context, vid uuid.UUID, oldLang, newLang string) error {
	query := `SELECT d.text FROM "dictionary" d 
				LEFT JOIN word w ON w.vocabulary_id=$1
				WHERE d.id=ANY(w.translate_ids);`

	rows, err := r.tr.Query(ctx, query, vid)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.MoveTranslatedWordsToNewDictionary: %w", err)
	}
	defer rows.Close()

	var texts []string
	for rows.Next() {
		var text string
		if err := rows.Scan(&text); err != nil {
			return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.MoveTranslatedWordsToNewDictionary: %w", err)
		}
		texts = append(texts, text)
	}

	if len(texts) == 0 {
		return nil
	}

	query = fmt.Sprintf(`INSERT INTO "%[2]s" SELECT 
		gen_random_uuid(), 
		text, 
		pronunciation, 
		'%[3]s' lang_code, 
		creator, 
		moderator, 
		updated_at, 
		created_at 
	FROM %[1]s 
	WHERE text=ANY($1::text[]) 
		AND lang_code=$2 ON CONFLICT DO NOTHING;`, getDictTable(oldLang), getDictTable(newLang), newLang)
	_, err = r.tr.Exec(ctx, query, texts, oldLang)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.MoveTranslatedWordsToNewDictionary: %w", err)
	}

	return nil
}

func (r *VocabRepo) UpdateVocabTranslatedIDs(ctx context.Context, vid uuid.UUID, oldLang, newLang string) error {
	query := `SELECT w.id, array_agg(d."text") FROM word w 
			LEFT JOIN "dictionary" d ON d.id = ANY(w.translate_ids)
			WHERE w.vocabulary_id=$1 
			AND array_length(w.translate_ids, 1) IS NOT NULL
			GROUP BY w.id;`

	rows, err := r.tr.Query(ctx, query, vid)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.UpdateVocabTranslatedIDs: %w", err)
	}
	defer rows.Close()

	type vocabWordData struct {
		ID    uuid.UUID
		Texts []string
	}

	var vocabWords []vocabWordData
	for rows.Next() {
		var vocabWord vocabWordData
		if err := rows.Scan(&vocabWord.ID, &vocabWord.Texts); err != nil {
			return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.UpdateVocabTranslatedIDs: %w", err)
		}
		vocabWords = append(vocabWords, vocabWord)
	}

	query = fmt.Sprintf(`UPDATE word SET 
		translate_ids=ARRAY(
			SELECT d.id FROM %s d 
			WHERE d.text=ANY($2::text[]) AND lang_code=$1) 
			WHERE id=$3;`, getDictTable(newLang))
	for _, vocabWord := range vocabWords {
		_, err := r.tr.Exec(ctx, query, newLang, vocabWord.Texts, vocabWord.ID)
		if err != nil {
			return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.UpdateVocabTranslatedIDs: %w", err)
		}
	}

	return nil
}

func (r *VocabRepo) DeleteTranslatedWordsFromOldDictionary(ctx context.Context, oldLang string) error {
	query := fmt.Sprintf(`DELETE FROM %[1]s WHERE id=ANY(SELECT d.id
		FROM %[1]s d
		LEFT JOIN word w ON d.id = ANY(w.translate_ids) OR d.id=w.native_id
		WHERE w.id IS NULL AND d.moderator IS NULL);`, getDictTable(oldLang))
	_, err := r.tr.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("vocabulary.delivery.repository.VocabRepo.DeleteTranslatedWordsFromOldDictionary: %w", err)
	}
	return nil
}
