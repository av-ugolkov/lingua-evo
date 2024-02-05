package repository

import (
	"context"
	"database/sql"
	"fmt"

	entity "github.com/av-ugolkov/lingua-evo/internal/services/language"
)

type LangRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *LangRepo {
	return &LangRepo{
		db: db,
	}
}

func (r *LangRepo) GetLanguage(ctx context.Context, langCode string) (*entity.Language, error) {
	query := `SELECT code, lang FROM language WHERE code=$1`
	language := entity.Language{}
	err := r.db.QueryRowContext(ctx, query, langCode).Scan(&language.Code, &language.Lang)
	if err != nil {
		return nil, fmt.Errorf("language.repository.LangRepo.GetLanguage - scan: %v", err)
	}

	return &language, nil
}

func (r *LangRepo) GetAvailableLanguages(ctx context.Context) ([]*entity.Language, error) {
	query := `SELECT code, lang FROM language ORDER BY lang ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("language.repository.LangRepo.GetAvailableLanguages: %v", err)
	}
	defer rows.Close()

	languages, err := scanRowsLanguage(rows)
	if err != nil {
		return nil, fmt.Errorf("language.repository.LangRepo.GetAvailableLanguages - scan: %v", err)
	}
	return languages, nil
}

func scanRowsLanguage(rows *sql.Rows) ([]*entity.Language, error) {
	var languages []*entity.Language
	for rows.Next() {
		var language entity.Language
		err := rows.Scan(
			&language.Code,
			&language.Lang,
		)
		if err != nil {
			return nil, err
		}

		languages = append(languages, &language)
	}

	return languages, nil
}
