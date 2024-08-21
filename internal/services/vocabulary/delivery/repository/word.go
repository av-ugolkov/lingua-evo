package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/vocabulary"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	// https://www.postgresql.org/docs/10/static/errcodes-appendix.html
	UniqueViolation = "23505"
)

func (r *VocabRepo) GetWord(ctx context.Context, id uuid.UUID) (entity.VocabWordData, error) {
	const query = `
	SELECT 
		w.id,
		w.vocabulary_id, 
		n.id native_id,
		n."text", 
		CASE WHEN w.pronunciation IS NOT NULL THEN w.pronunciation ELSE '' END pronunciation, 
		n.lang_code, 
		array_agg(distinct t."text") FILTER (WHERE t."text" IS NOT NULL) translates, 
		array_agg(distinct e."text") FILTER (WHERE e."text" IS NOT NULL) examples,
		w.updated_at,
		w.created_at 
	FROM word w
		LEFT JOIN "dictionary" n ON n.id = w.native_id
		LEFT JOIN "dictionary" t ON t.id = ANY(w.translate_ids)
		LEFT JOIN "example" e ON e.id = ANY(w.example_ids)
	WHERE w.id=$1
	GROUP BY w.id, n.id, n."text", w.pronunciation, n.lang_code;`

	var vocabWordData entity.VocabWordData
	var translates []string
	var examples []string
	err := r.pgxPool.QueryRow(ctx, query, id).Scan(
		&vocabWordData.ID,
		&vocabWordData.VocabID,
		&vocabWordData.Native.ID,
		&vocabWordData.Native.Text,
		&vocabWordData.Native.Pronunciation,
		&vocabWordData.Native.LangCode,
		&translates,
		&examples,
		&vocabWordData.UpdatedAt,
		&vocabWordData.CreatedAt)
	if err != nil {
		return entity.VocabWordData{}, fmt.Errorf("word.repository.WordRepo.GetWord: %w", err)
	}

	for _, tr := range translates {
		vocabWordData.Translates = append(vocabWordData.Translates, entityDict.DictWord{Text: tr})
	}

	for _, ex := range examples {
		vocabWordData.Examples = append(vocabWordData.Examples, entityExample.Example{Text: ex})
	}

	return vocabWordData, nil
}

func (r *VocabRepo) AddWord(ctx context.Context, word entity.VocabWord) (uuid.UUID, error) {
	const query = `
	INSERT INTO word (
		id,
		vocabulary_id,
		native_id, 
		pronunciation, 
		translate_ids, 
		example_ids, 
		updated_at, 
		created_at) 
	VALUES($1, $2, $3, $4, $5, $6, $7, $7);`
	vocabWordID := uuid.New()
	_, err := r.pgxPool.Exec(ctx, query, vocabWordID, word.VocabID, word.NativeID, word.Pronunciation, word.TranslateIDs, word.ExampleIDs, time.Now().UTC())
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == UniqueViolation:
			return uuid.Nil, fmt.Errorf("word.repository.WordRepo.AddWord: %w", entity.ErrDuplicate)
		default:
			return uuid.Nil, fmt.Errorf("word.repository.WordRepo.AddWord: %w", err)
		}
	}

	return vocabWordID, nil
}

func (r *VocabRepo) GetWordsFromVocabulary(ctx context.Context, dictID uuid.UUID, capacity int) ([]string, error) {
	const query = `
	SELECT text 
	FROM dictionary 
	WHERE id=any(
		SELECT native_id
		FROM word 
		WHERE vocabulary_id=$1
			ORDER BY random() LIMIT $2)`
	rows, err := r.pgxPool.Query(ctx, query, dictID, capacity)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetWordsFromVocabulary: %w", err)
	}
	defer rows.Close()

	words := make([]string, 0, capacity)
	for rows.Next() {
		var word string
		err = rows.Scan(&word)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetWordsFromVocabulary - scan: %w", err)
		}
		words = append(words, word)
	}

	return words, nil
}

func (r *VocabRepo) GetRandomWord(ctx context.Context, vocabID uuid.UUID) (entity.VocabWord, error) {
	var vocabWord entity.VocabWord
	query := `SELECT native_id, translate_ids, example_ids FROM word WHERE vocabulary_id=$1 ORDER BY random() LIMIT 1;`
	err := r.pgxPool.QueryRow(ctx, query, vocabID).Scan(
		&vocabWord.NativeID,
		&vocabWord.TranslateIDs,
		&vocabWord.ExampleIDs,
	)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.repository.WordRepo.GetRandomWord - scan: %w", err)
	}
	return vocabWord, nil
}

func (r *VocabRepo) DeleteWord(ctx context.Context, vocabWord entity.VocabWord) error {
	query := `DELETE FROM word WHERE vocabulary_id=$1 AND id=$2;`
	result, err := r.pgxPool.Exec(ctx, query, vocabWord.VocabID, vocabWord.ID)
	if err != nil {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord - exec: %w", err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord - rows affected: %w", pgx.ErrNoRows)
	}

	return nil
}

func (r *VocabRepo) GetRandomVocabulary(ctx context.Context, vocabID uuid.UUID, limit int) ([]entity.VocabWordData, error) {
	query := `
	SELECT 
		n.id native_id,
		n."text", 
		CASE WHEN w.pronunciation IS NOT NULL THEN w.pronunciation ELSE '' END pronunciation, 
		array_agg(distinct t."text") FILTER (WHERE t."text" IS NOT NULL) translates
	FROM word w
		LEFT JOIN "dictionary" n ON n.id = w.native_id
		LEFT JOIN "dictionary" t ON t.id = ANY(w.translate_ids)
	WHERE vocabulary_id=$1
	GROUP BY n.id, n."text", w.pronunciation
	ORDER BY RANDOM() 
	LIMIT $2;`
	rows, err := r.pgxPool.Query(ctx, query, vocabID, limit)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetRandomVocabulary: %w", err)
	}
	defer rows.Close()

	vocabularyWords := make([]entity.VocabWordData, 0, limit)
	for rows.Next() {
		var vocabulary entity.VocabWordData
		var translates []string
		err = rows.Scan(
			&vocabulary.Native.ID,
			&vocabulary.Native.Text,
			&vocabulary.Native.Pronunciation,
			&translates,
		)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetRandomVocabulary - scan: %w", err)
		}

		for _, tr := range translates {
			vocabulary.Translates = append(vocabulary.Translates, entityDict.DictWord{Text: tr})
		}

		vocabularyWords = append(vocabularyWords, vocabulary)
	}

	return vocabularyWords, nil
}

func (r *VocabRepo) GetVocabulary(ctx context.Context, vocabID uuid.UUID) ([]entity.VocabWord, error) {
	query := `SELECT id, native_id, translate_ids, example_ids, updated_at, created_at FROM word WHERE vocabulary_id=$1;`
	rows, err := r.pgxPool.Query(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabulary: %w", err)
	}
	defer rows.Close()

	vocabularies := make([]entity.VocabWord, 0, 25)
	for rows.Next() {
		var vocabulary entity.VocabWord
		err = rows.Scan(
			&vocabulary.ID,
			&vocabulary.NativeID,
			&vocabulary.TranslateIDs,
			&vocabulary.ExampleIDs,
			&vocabulary.UpdatedAt,
			&vocabulary.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetVocabulary - scan: %w", err)
		}
		vocabularies = append(vocabularies, vocabulary)
	}

	return vocabularies, nil
}

func (r *VocabRepo) GetVocabularyWords(ctx context.Context, vocabID uuid.UUID) ([]entity.VocabWordData, error) {
	var countRows int
	err := r.pgxPool.QueryRow(ctx, `SELECT count(*) FROM word WHERE vocabulary_id=$1`, vocabID).Scan(&countRows)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords - count: %w", err)
	}

	var langNative, langTranslate string
	err = r.pgxPool.QueryRow(ctx, `SELECT native_lang, translate_lang FROM vocabulary WHERE id=$1`, vocabID).Scan(&langNative, &langTranslate)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords - get langs: %w", err)
	}

	query := fmt.Sprintf(`
	SELECT
		w.id,
		n.id native_id,
		n."text",
		CASE WHEN w.pronunciation IS NOT NULL THEN w.pronunciation ELSE '' END pronunciation,
		array_agg(distinct t."text") FILTER (WHERE t."text" IS NOT NULL) translates,
		array_agg(distinct e."text") FILTER (WHERE e."text" IS NOT NULL) examples,
		w.updated_at,
		w.created_at
	FROM word w
		LEFT JOIN "dictionary_%[1]s" n ON n.id = w.native_id
		LEFT JOIN "dictionary_%[2]s" t ON t.id = ANY(w.translate_ids)
		LEFT JOIN "example_%[1]s" e ON e.id = ANY(w.example_ids)
	WHERE w.vocabulary_id=$1
	GROUP BY w.id, n.id, n."text", w.pronunciation, n.lang_code
	LIMIT $2;`, langNative, langTranslate)

	rows, err := r.pgxPool.Query(ctx, query, vocabID, countRows)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords: %w", err)
	}
	defer rows.Close()

	vocabularyWords := make([]entity.VocabWordData, 0, countRows)
	for rows.Next() {
		var wordData entity.VocabWordData
		var translates []string
		var examples []string
		err = rows.Scan(
			&wordData.ID,
			&wordData.Native.ID,
			&wordData.Native.Text,
			&wordData.Native.Pronunciation,
			&translates,
			&examples,
			&wordData.UpdatedAt,
			&wordData.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords - scan: %w", err)
		}

		for _, tr := range translates {
			wordData.Translates = append(wordData.Translates, entityDict.DictWord{Text: tr, LangCode: langTranslate})
		}

		for _, ex := range examples {
			wordData.Examples = append(wordData.Examples, entityExample.Example{Text: ex})
		}

		wordData.VocabID = vocabID
		wordData.Native.LangCode = langNative

		vocabularyWords = append(vocabularyWords, wordData)
	}

	return vocabularyWords, nil
}

func (r *VocabRepo) UpdateWord(ctx context.Context, vocabWord entity.VocabWord) error {
	query := `UPDATE word SET native_id=$1, pronunciation=$2, translate_ids=$3, example_ids=$4, updated_at=$5 WHERE id=$6;`

	result, err := r.pgxPool.Exec(ctx, query, vocabWord.NativeID, vocabWord.Pronunciation, vocabWord.TranslateIDs, vocabWord.ExampleIDs, vocabWord.UpdatedAt.Format(time.RFC3339), vocabWord.ID)
	if err != nil {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord - exec: %w", err)
	}

	if rows := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord - rows affected: %w", pgx.ErrNoRows)
	}

	return nil
}

func (r *VocabRepo) GetCountWords(ctx context.Context, userID uuid.UUID) (int, error) {
	const query = `SELECT count(id) FROM word WHERE vocabulary_id=ANY(SELECT id FROM vocabulary WHERE user_id=$1);`

	var count int
	err := r.pgxPool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("word.repository.WordRepo.GetCountWords: %w", err)
	}

	return count, nil
}
