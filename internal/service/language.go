package service

import (
	"context"
	"fmt"
)

type Language struct {
	Code string
	Lang string
}

func (l *Lingua) GetLanguages(ctx context.Context) ([]*Language, error) {
	dbLanguages, err := l.db.GetLanguages(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.lingua.GetLanguages: %v", err)
	}
	languages := make([]*Language, 0, len(dbLanguages))
	for _, language := range dbLanguages {
		languages = append(languages, &Language{language.Code, language.Lang})
	}

	return languages, nil
}
