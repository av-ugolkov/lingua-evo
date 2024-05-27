package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	entityDict "github.com/av-ugolkov/lingua-evo/internal/services/dictionary"
	entityExample "github.com/av-ugolkov/lingua-evo/internal/services/example"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/word"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
)

type WordRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *WordRepo {
	return &WordRepo{
		db: db,
	}
}

func (r *WordRepo) GetWord(ctx context.Context, id uuid.UUID) (entity.VocabWordData, error) {
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
	WHERE id=$1
	GROUP BY w.id, n.id, n."text", w.pronunciation, n.lang_code;`

	var vocabWordData entity.VocabWordData
	var translates []string
	var examples []string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vocabWordData.ID,
		&vocabWordData.Native.ID,
		&vocabWordData.Native.Text,
		&vocabWordData.Native.Pronunciation,
		&vocabWordData.Native.LangCode,
		pq.Array(&translates),
		pq.Array(&examples),
		&vocabWordData.UpdatedAt,
		&vocabWordData.CreatedAt)
	if err != nil {
		return entity.VocabWordData{}, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords: %w", err)
	}

	for _, tr := range translates {
		vocabWordData.Translates = append(vocabWordData.Translates, entityDict.DictWord{Text: tr})
	}

	for _, ex := range examples {
		vocabWordData.Examples = append(vocabWordData.Examples, entityExample.Example{Text: ex})
	}

	return vocabWordData, nil
}

func (r *WordRepo) AddWord(ctx context.Context, word entity.VocabWord) error {
	const query = `INSERT INTO word (id, vocabulary_id, native_id, translate_ids, example_ids, updated_at, created_at) VALUES($1, $2, $3, $4, $5, $6, $6);`
	_, err := r.db.ExecContext(ctx, query, word.ID, word.VocabID, word.NativeID, word.TranslateIDs, word.ExampleIDs, time.Now().UTC())
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == "23505":
			return fmt.Errorf("word.repository.WordRepo.AddWord: %w", entity.ErrDuplicate)
		default:
			return fmt.Errorf("word.repository.WordRepo.AddWord: %w", err)
		}
	}

	return nil
}

func (r *WordRepo) GetWordsFromVocabulary(ctx context.Context, dictID uuid.UUID, capacity int) ([]string, error) {
	const query = `
	SELECT text 
	FROM dictionary 
	WHERE id=any(
		SELECT native_id
		FROM word 
		WHERE vocabulary_id=$1
			ORDER BY random() LIMIT $2)`
	rows, err := r.db.QueryContext(ctx, query, dictID, capacity)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetWordsFromVocabulary: %w", err)
	}
	defer func() { _ = rows.Close() }()

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

func (r *WordRepo) GetRandomWord(ctx context.Context, vocabID uuid.UUID) (entity.VocabWord, error) {
	var vocabWord entity.VocabWord
	query := `SELECT native_id, translate_ids, example_ids FROM word WHERE vocabulary_id=$1 ORDER BY random() LIMIT 1;`
	err := r.db.QueryRowContext(ctx, query, vocabID).Scan(
		&vocabWord.NativeID,
		&vocabWord.TranslateIDs,
		&vocabWord.ExampleIDs,
	)
	if err != nil {
		return entity.VocabWord{}, fmt.Errorf("word.repository.WordRepo.GetRandomWord - scan: %w", err)
	}
	return vocabWord, nil
}

func (r *WordRepo) DeleteWord(ctx context.Context, vocabWord entity.VocabWord) error {
	query := `DELETE FROM word WHERE vocabulary_id=$1 AND id=$2;`
	result, err := r.db.ExecContext(ctx, query, vocabWord.VocabID, vocabWord.ID)
	if err != nil {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord - exec: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord - rows affected: %w", sql.ErrNoRows)
	} else if err != nil {
		return fmt.Errorf("word.repository.WordRepo.DeleteWord: %w", err)
	}
	return nil
}

func (r *WordRepo) GetRandomVocabulary(ctx context.Context, vocabID uuid.UUID, limit int) ([]entity.VocabWordData, error) {
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
	rows, err := r.db.QueryContext(ctx, query, vocabID, limit)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetRandomVocabulary: %w", err)
	}
	defer func() { _ = rows.Close() }()

	vocabularyWords := make([]entity.VocabWordData, 0, limit)
	for rows.Next() {
		var vocabulary entity.VocabWordData
		var translates []string
		err = rows.Scan(
			&vocabulary.Native.ID,
			&vocabulary.Native.Text,
			&vocabulary.Native.Pronunciation,
			pq.Array(&translates),
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

func (r *WordRepo) GetVocabulary(ctx context.Context, vocabID uuid.UUID) ([]entity.VocabWord, error) {
	query := `SELECT id, native_id, translate_ids, example_ids, updated_at, created_at FROM word WHERE vocabulary_id=$1;`
	rows, err := r.db.QueryContext(ctx, query, vocabID)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabulary: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	vocabularies := make([]entity.VocabWord, 0, 25)
	for rows.Next() {
		var vocabulary entity.VocabWord
		err = rows.Scan(
			&vocabulary.ID,
			&vocabulary.NativeID,
			pq.Array(&vocabulary.TranslateIDs),
			pq.Array(&vocabulary.ExampleIDs),
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

func (r *WordRepo) GetVocabularyWords(ctx context.Context, vocabID uuid.UUID) ([]entity.VocabWordData, error) {
	var countRows int
	err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM word WHERE vocabulary_id=$1`, vocabID).Scan(&countRows)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords - count: %w", err)
	}

	var langNative, langTranslate string
	err = r.db.QueryRowContext(ctx, `SELECT native_lang, translate_lang FROM vocabulary WHERE id=$1`, vocabID).Scan(&langNative, &langTranslate)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords - get langs: %w", err)
	}

	query := fmt.Sprintf(`
	SELECT 
		w.id, 
		n.id native_id,
		n."text", 
		CASE WHEN w.pronunciation IS NOT NULL THEN w.pronunciation ELSE '' END pronunciation, 
		n.lang_code, 
		array_agg(distinct t."text") FILTER (WHERE t."text" IS NOT NULL) translates, 
		array_agg(distinct e."text") FILTER (WHERE e."text" IS NOT NULL) examples,
		w.updated_at,
		w.created_at 
	FROM word w
		LEFT JOIN "dictionary_%[1]s" n ON n.id = w.native_id
		LEFT JOIN "dictionary_%[2]s" t ON t.id = ANY(w.translate_ids)
		LEFT JOIN "example_%[1]s" e ON e.id = ANY(w.example_ids)
	WHERE vocabulary_id=$1
	GROUP BY w.id, n.id, n."text", w.pronunciation, n.lang_code
	LIMIT $2;`, langNative, langTranslate)

	rows, err := r.db.QueryContext(ctx, query, vocabID, countRows)
	if err != nil {
		return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	vocabularyWords := make([]entity.VocabWordData, 0, countRows)
	for rows.Next() {
		var vocabulary entity.VocabWordData
		var translates []string
		var examples []string
		err = rows.Scan(
			&vocabulary.ID,
			&vocabulary.Native.ID,
			&vocabulary.Native.Text,
			&vocabulary.Native.Pronunciation,
			&vocabulary.Native.LangCode,
			pq.Array(&translates),
			pq.Array(&examples),
			&vocabulary.UpdatedAt,
			&vocabulary.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("word.repository.WordRepo.GetVocabularyWords - scan: %w", err)
		}

		for _, tr := range translates {
			vocabulary.Translates = append(vocabulary.Translates, entityDict.DictWord{Text: tr})
		}

		for _, ex := range examples {
			vocabulary.Examples = append(vocabulary.Examples, entityExample.Example{Text: ex})
		}

		vocabularyWords = append(vocabularyWords, vocabulary)
	}

	return vocabularyWords, nil
}

func (r *WordRepo) UpdateWord(ctx context.Context, vocabWord entity.VocabWord) error {
	query := `UPDATE word SET native_id=$1, pronunciation=$2, translate_ids=$3, example_ids=$4, updated_at=$5 WHERE id=$6;`

	result, err := r.db.ExecContext(ctx, query, vocabWord.NativeID, vocabWord.Pronunciation, vocabWord.TranslateIDs, vocabWord.ExampleIDs, vocabWord.UpdatedAt.Format(time.RFC3339), vocabWord.ID)
	if err != nil {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord - exec: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord - rows affected: %w", sql.ErrNoRows)
	} else if err != nil {
		return fmt.Errorf("word.repository.WordRepo.UpdateWord: %w", err)
	}
	return nil
}

func (r *WordRepo) GetCountWords(ctx context.Context, userID uuid.UUID) (int, error) {
	const query = `SELECT count(id) FROM word WHERE vocabulary_id=ANY(SELECT id FROM vocabulary WHERE user_id=$1);`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("word.repository.WordRepo.GetCountWords: %w", err)
	}

	return count, nil
}
