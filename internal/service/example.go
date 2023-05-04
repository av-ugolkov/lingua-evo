package service

import (
	"context"

	"github.com/google/uuid"
)

func (l *Lingua) AddExample(ctx context.Context, wordId uuid.UUID, example string) (uuid.UUID, error) {
	exampleId, err := l.db.AddExample(ctx, wordId, example)
	if err != nil {
		return uuid.Nil, err
	}

	return exampleId, nil
}
