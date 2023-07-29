package services

import (
	"context"

	"github.com/google/uuid"
)

func (l *Lingua) AddWordInDictionary(ctx context.Context, userID, origWordId uuid.UUID, tranWordId []uuid.UUID, pronunciation string, examples []uuid.UUID) (uuid.UUID, error) {
	err := l.db.AddWordInDictionary(ctx, userID, origWordId, tranWordId, pronunciation, examples)
	if err != nil {
		return uuid.Nil, err
	}

	return uuid.Nil, nil
}
