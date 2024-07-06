package repository

import (
	"context"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/language"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LangRepo struct {
	pgxPool *pgxpool.Pool
}

func NewRepo(pgxPool *pgxpool.Pool) *LangRepo {
	return &LangRepo{
		pgxPool: pgxPool,
	}
}

func (r *LangRepo) GetLanguage(ctx context.Context, langCode string) (string, error) {
	query := `SELECT lang FROM language WHERE code=$1`
	var language string
	err := r.pgxPool.QueryRow(ctx, query, langCode).Scan(&language)
	if err != nil {
		return "", fmt.Errorf("language.repository.LangRepo.GetLanguage: %v", err)
	}

	return language, nil
}

func (r *LangRepo) GetAvailableLanguages(ctx context.Context) ([]entity.Language, error) {
	query := `SELECT code, lang FROM language ORDER BY lang`
	rows, err := r.pgxPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("language.repository.LangRepo.GetAvailableLanguages: %v", err)
	}
	defer rows.Close()

	languages := make([]entity.Language, 0)
	for rows.Next() {
		var lang entity.Language

		err := rows.Scan(&lang.Code, &lang.Lang)
		if err != nil {
			return nil, fmt.Errorf("language.repository.LangRepo.GetAvailableLanguages: %v", err)
		}

		languages = append(languages, lang)
	}

	return languages, nil
}
