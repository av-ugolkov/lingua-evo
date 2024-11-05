package repository

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/db/transactor"
	entity "github.com/av-ugolkov/lingua-evo/internal/services/language"
	"github.com/av-ugolkov/lingua-evo/runtime"
)

type LangRepo struct {
	tr *transactor.Transactor
}

func NewRepo(tr *transactor.Transactor) *LangRepo {
	return &LangRepo{
		tr: tr,
	}
}

func (r *LangRepo) GetLanguage(ctx context.Context, langCode string) (string, error) {
	query := `SELECT lang FROM language WHERE code=$1`
	var language string
	err := r.tr.QueryRow(ctx, query, langCode).Scan(&language)
	if err != nil {
		return runtime.EmptyString, fmt.Errorf("language.repository.LangRepo.GetLanguage: %v", err)
	}

	return language, nil
}

func (r *LangRepo) GetAvailableLanguages(ctx context.Context) ([]entity.Language, error) {
	query := `SELECT code, lang FROM language ORDER BY lang`
	rows, err := r.tr.Query(ctx, query)
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
