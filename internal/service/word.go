package service

import (
	"context"

	"github.com/google/uuid"

	"lingua-evo/internal/delivery/repository"
)

func (l *Lingua) AddWord(ctx context.Context, word *repository.Word) (uuid.UUID, error) {
	wordId, err := l.db.FindWord(ctx, word)
	if err == nil {
		return wordId, nil
	}
	l.logger.Warn(err)

	wordId, err = l.db.AddWord(ctx, word)
	if err != nil {
		return uuid.UUID{}, err
	}

	return wordId, nil
}

func (l *Lingua) EditWord(ctx context.Context, w *repository.Word) error {
	return nil
}

func (l *Lingua) FindWord(ctx context.Context, w *repository.Word) (*repository.Word, error) {
	return nil, nil
}

func (l *Lingua) FindWords(ctx context.Context, w string) (*repository.Word, error) {
	return nil, nil
}

func (l *Lingua) RemoveWord(ctx context.Context, w *repository.Word) error {
	return nil
}

func (l *Lingua) PickRandomWord(ctx context.Context, w *repository.Word) (*repository.Word, error) {
	return nil, nil
}

func (l *Lingua) SharedWord(ctx context.Context, w *repository.Word) (*repository.Word, error) {
	return nil, nil
}
