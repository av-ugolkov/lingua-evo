package repository

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/games"

	"github.com/google/uuid"
)

type Repo struct {
	tr *transactor.Transactor
}

func New(tr *transactor.Transactor) *Repo {
	return &Repo{
		tr: tr,
	}
}

func (r *Repo) GerWords(ctx context.Context, uid, vid uuid.UUID, count int) ([]entity.ReviseGameWord, error) {
	query := `
		SELECT 
			dn.text, 
			array_agg(dt.text), 
			array_agg(en.text), 
			COALESCE(grs.right, 0) grsr, 
			COALESCE(grs.wrong, 0) grsw 
		FROM word w
		LEFT JOIN game_revise_stats grs ON grs.vocab_word_id = w.id AND grs.user_id = $1
		LEFT JOIN dictionary dn ON dn.id = w.native_id
		LEFT JOIN dictionary dt ON dt.id = ANY(w.translate_ids)
		LEFT JOIN example en ON en.id = ANY(w.example_ids) 
		WHERE w.vocabulary_id = $2
		GROUP BY dn.text, w.example_ids, grsr, grsw
		ORDER BY grsr, grsw
		LIMIT $3;`

	rows, err := r.tr.Query(ctx, query, uid, vid, count)
	if err != nil {
		return nil, fmt.Errorf("subscribers.repository.Repo.Get: %w", err)
	}
	defer rows.Close()

	words := make([]entity.ReviseGameWord, 0, count)
	for rows.Next() {
		var word entity.ReviseGameWord
		err := rows.Scan(
			&word.Text,
			&word.Translates,
			&word.Examples,
			&word.Right,
			&word.Wrong)
		if err != nil {
			return nil, fmt.Errorf("subscribers.repository.Repo.Get: %w", err)
		}

		words = append(words, word)
	}

	return words, nil
}
